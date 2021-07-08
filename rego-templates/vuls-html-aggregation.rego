package postee.vuls.html.aggregation

import data.postee.flat_array


title := "Vulnerability scan report"
result := res {
    scans := [ scan | 
            item:=input[i].description

            scan:=[sprintf("<h1>%s</h1>", [input[i].title]), item]
    ] 

    res:= concat("\n", flat_array(scans))
}


