func main() {
    input := []string{"flower", "flow", "flight"}
    fmt.Println(longestCommonPrefix(input)) // "fl"
    
    input2 := []string{"dog", "racecar", "car"}
    fmt.Println(longestCommonPrefix(input2)) // ""
}



package main

import (
	"strings"
)

func longestCommonPrefix(strs []string) string {
    if len(strs) == 0 {
        return ""
    }

    prefix := strs[0]

    for i := 1; i < len(strs); i++ {
        for !strings.HasPrefix(strs[i], prefix) {
            prefix = prefix[:len(prefix)-1]
            if prefix == "" {              
                return ""
            }
        }
    }
    return prefix
}


//    0 1 2 3 4 5 
//----------------
// 0 |f l o w e r
// 2 |f l o w
// 3 |f l i g h t
// 