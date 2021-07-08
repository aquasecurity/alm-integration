class validator {
    constructor(fields, validationFn) {
        this.validationFn = validationFn
        this.register = (id, label, name) => {
            this.label = label
            this.name = name
            fields[id] = this
        }
    }
}
export default {
    methods:
    {
        url(label, value) {
            if (!value) {
                return false
            }
            const errorMsg = `Invalid url : ${value}`
            let url

            try {
                url = new URL(value);
            } catch (_) {
                return errorMsg;
            }

            return url.protocol === "http:" || url.protocol === "https:" ? false : errorMsg;
        },

        email(label, value) {
            const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
            return re.test(String(value).toLowerCase()) ? false : `Invalid email '${value}'`
        },

        required(label, value) {
            return !value ? `${label} is required` : false;
        },

        recipients(label, value) {
            const hasOneElement = value && value.length && value[0]
            if (!hasOneElement) {
                return `At least one of ${label} is required`
            } else {
                for (const email of value) {
                    const v = this.email("-", email);
                    if (v) {
                        return v;
                    }
                }
            }
            return false
        },
        v(validationFn) {
            return new validator(this.fields, validationFn);
        }


    }
}