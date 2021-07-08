<template>
<div >
    <div class="row justify-content-end pb-3 pr-3">
        <router-link :to="{name: 'add-template'}" class="btn btn-primary">Add Template</router-link>
    </div>
    <div class="row row-cols-1 row-cols-md-3">
        <TemplateCard v-for="(template, index) in templates" :key="'template-' + index" :name="template.name"/>
    </div>
</div>
</template>
<script>
import TemplateCard from './TemplateCard.vue';
import { mapState} from 'vuex'

export default {
    components: {
        TemplateCard
    },
    computed: {
        ...mapState({
            templates (state) {
                return state.templates.all
            }
        })
    }

}
</script>