package postee
############################################# Common functions ############################################
by_flag(a, b, flag) = a {
	flag
}
by_flag(a, b, flag) = b {
	flag = false
}
duplicate(a, b, col) = a {col == 1}
duplicate(a, b, col) = b {col == 2}

clamp(a, b) = b { a > b }
clamp(a, b) = a { a <= b }

by_flag(a, b, flag) = a {
	flag
}
by_flag(a, b, flag) = b {
	flag = false
}
flat_array(a) = o {
	o:=[item |
    	item:=a[_][_]
    ]
}
with_default(obj, prop, default_value) = default_value{
 not obj[prop]
}
with_default(obj, prop, default_value) = obj[prop]{
 obj[prop]
}
