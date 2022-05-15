package releaseupdater

import (
	"fmt"
	"regexp"
)

/*
returns a copy of content where versions according to the regex are replaced.

Example:
 (?P<key>\s*version:\s*)(?P<value>v.*?)(?P<other>$|\s+\.*)

Important: the regex MUST define the following groups keys and other in order to work.
Otherwise the version is not replaced.

key defines the key including the preceeding characters. This may be needed for
formats like YAML.

other contains everything you want to keep in that line behind the version value to replace.
*/
func ReplaceVersion(content []byte, regex, newVersion string) []byte {

	r, err := regexp.Compile(regex)
	if err != nil {
		fmt.Printf("Error compiling regex %v: %v\n", regex, err)
		return content
	}

	ret := r.ReplaceAll(content, []byte(fmt.Sprintf("${key}%s${other}", newVersion)))
	// fmt.Println(string(ret))
	return ret
}

// convenience method for ReplaceVersion to ease testing
func ReplaceVersionString(content, regex, newVersion string) string {
	return string(ReplaceVersion([]byte(content), regex, newVersion))
}
