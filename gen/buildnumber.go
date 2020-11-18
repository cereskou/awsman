// +build ignore

//input file: build.json versioninfo.json
//output    : version.go created.txt
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

//VersionInfo -
type VersionInfo struct {
	FixedFileInfo struct {
		FileVersion struct {
			Major int `json:"Major"`
			Minor int `json:"Minor"`
			Patch int `json:"Patch"`
			Build int `json:"Build"`
		} `json:"FileVersion"`
		ProductVersion struct {
			Major int `json:"Major"`
			Minor int `json:"Minor"`
			Patch int `json:"Patch"`
			Build int `json:"Build"`
		} `json:"ProductVersion"`
		FileFlagsMask string `json:"FileFlagsMask"`
		FileFlags     string `json:"FileFlags "`
		FileOS        string `json:"FileOS"`
		FileType      string `json:"FileType"`
		FileSubType   string `json:"FileSubType"`
	} `json:"FixedFileInfo"`
	StringFileInfo struct {
		Comments         string `json:"Comments"`
		CompanyName      string `json:"CompanyName"`
		FileDescription  string `json:"FileDescription"`
		FileVersion      string `json:"FileVersion"`
		InternalName     string `json:"InternalName"`
		LegalCopyright   string `json:"LegalCopyright"`
		LegalTrademarks  string `json:"LegalTrademarks"`
		OriginalFilename string `json:"OriginalFilename"`
		PrivateBuild     string `json:"PrivateBuild"`
		ProductName      string `json:"ProductName"`
		ProductVersion   string `json:"ProductVersion"`
		SpecialBuild     string `json:"SpecialBuild"`
	} `json:"StringFileInfo"`
	VarFileInfo struct {
		Translation struct {
			LangID    string `json:"LangID"`
			CharsetID string `json:"CharsetID"`
		} `json:"Translation"`
	} `json:"VarFileInfo"`
	IconPath     string `json:"IconPath"`
	ManifestPath string `json:"ManifestPath"`
}

//BuildVersion -
type BuildVersion struct {
	Major            int    `json:"Major"`
	Minor            int    `json:"Minor"`
	Patch            int    `json:"Patch"`
	Build            int    `json:"Build"`
	Date             string `json:"Date"`
	Time             string `json:"Time"`
	CompanyName      string `json:"CompanyName"`
	FileDescription  string `json:"FileDescription"`
	InternalName     string `json:"InternalName"`
	LegalCopyright   string `json:"LegalCopyright"`
	OriginalFilename string `json:"OriginalFilename"`
	ProductName      string `json:"ProductName"`
	IconPath         string `json:"IconPath"`
	ManifestPath     string `json:"ManifestPath"`
}

//version
var (
	buildver     string = "build.json"
	fversion     string = "version.go"
	fversioninfo string = "versioninfo.json"
	fbuild       string = "created.txt"
)

//var default versioninfo
var versioninfo string = `
{
    "FixedFileInfo": {
        "FileVersion": {
            "Major": 0,
            "Minor": 0,
            "Patch": 0,
            "Build": 0
        },
        "ProductVersion": {
            "Major": 0,
            "Minor": 0,
            "Patch": 0,
            "Build": 0
        },
        "FileFlagsMask": "3f",
        "FileFlags ": "00",
        "FileOS": "040004",
        "FileType": "01",
        "FileSubType": "00"
    },
    "StringFileInfo": {
        "Comments": "",
        "CompanyName": "",
        "FileDescription": "",
        "FileVersion": "",
        "InternalName": "",
        "LegalCopyright": "",
        "LegalTrademarks": "",
        "OriginalFilename": "",
        "PrivateBuild": "",
        "ProductName": "",
        "ProductVersion": "",
        "SpecialBuild": ""
    },
    "VarFileInfo": {
        "Translation": {
            "LangID": "0411",
            "CharsetID": "04B0"
        }
    },
    "IconPath": "",
    "ManifestPath": ""
}`

func main() {
	b, err := ioutil.ReadFile(buildver)
	if err != nil {
		log.Printf("Cannot open %q: %v\n", buildver, err)
		os.Exit(-1)
	}
	var ver BuildVersion
	err = json.Unmarshal(b, &ver)
	if err != nil {
		log.Printf("Cannot decode json: %v\n", err)
		os.Exit(-1)
	}

	b, err = ioutil.ReadFile(fversioninfo)
	if err != nil {
		//log.Printf("Cannot open %q: %v\n", fversioninfo, err)
		b = []byte(versioninfo)
	}

	var vi VersionInfo
	err = json.Unmarshal(b, &vi)
	if err != nil {
		log.Printf("Cannot decode json: %v\n", err)
		os.Exit(-1)
	}

	//build number plus
	ver.Build++

	//update build.json
	jsontxt, _ := json.MarshalIndent(ver, "", "    ")
	err = ioutil.WriteFile(buildver, jsontxt, 0644)
	if err != nil {
		log.Printf("Cannot write json: %v\n", err)
		os.Exit(-1)
	}

	err = writeVersion(&ver)
	if err != nil {
		log.Printf("Cannot write to version.go: %v\n", err)
		os.Exit(-1)

	}

	if len(ver.Date) > 0 && len(ver.Time) > 0 {
		ftime := fmt.Sprintf("%v %v", ver.Date, ver.Time)
		err = ioutil.WriteFile(fbuild, []byte(ftime), 0644)
		if err != nil {
			log.Printf("Cannot write %v: %v\n", ftime, err)
			os.Exit(-1)
		}
	}

	//versioninfo - FileVersion
	vi.FixedFileInfo.FileVersion.Major = ver.Major
	vi.FixedFileInfo.FileVersion.Minor = ver.Minor
	vi.FixedFileInfo.FileVersion.Patch = ver.Patch
	vi.FixedFileInfo.FileVersion.Build = ver.Build
	//versioninfo - ProductVersion
	vi.FixedFileInfo.ProductVersion.Major = ver.Major
	vi.FixedFileInfo.ProductVersion.Minor = ver.Minor
	vi.FixedFileInfo.ProductVersion.Patch = ver.Patch
	vi.FixedFileInfo.ProductVersion.Build = ver.Build

	fver := fmt.Sprintf("%v.%v.%v.%v", ver.Major, ver.Minor, ver.Patch, ver.Build)
	//FileVersion
	vi.StringFileInfo.FileVersion = fver
	pver := fmt.Sprintf("%v.%v.%v", ver.Major, ver.Minor, ver.Patch)
	//ProductVersion
	vi.StringFileInfo.ProductVersion = pver

	if len(ver.CompanyName) > 0 {
		vi.StringFileInfo.CompanyName = ver.CompanyName
	}
	if len(ver.FileDescription) > 0 {
		vi.StringFileInfo.FileDescription = ver.FileDescription
	}
	if len(ver.LegalCopyright) > 0 {
		vi.StringFileInfo.LegalCopyright = ver.LegalCopyright
	}
	if len(ver.ProductName) > 0 {
		vi.StringFileInfo.ProductName = ver.ProductName
		vi.StringFileInfo.InternalName = ver.ProductName
		vi.StringFileInfo.OriginalFilename = ver.ProductName + ".exe"
	}
	if len(ver.InternalName) > 0 {
		vi.StringFileInfo.InternalName = ver.InternalName
	}
	if len(ver.OriginalFilename) > 0 {
		vi.StringFileInfo.OriginalFilename = ver.OriginalFilename
	}
	if len(ver.IconPath) > 0 {
		vi.IconPath = ver.IconPath
	}
	if len(ver.ManifestPath) > 0 {
		vi.ManifestPath = ver.ManifestPath
	}

	//update build.json
	jsontxt, _ = json.MarshalIndent(vi, "", "    ")
	err = ioutil.WriteFile(fversioninfo, jsontxt, 0644)
	if err != nil {
		log.Printf("Cannot write json: %v\n", err)
		os.Exit(-1)
	}

	os.Exit(0)
}

func writeVersion(v *BuildVersion) error {
	var tmpl = `package main

//var -
var (
	version string = %q
)
`
	date := strings.ReplaceAll(v.Date, "/", "")
	// ver := fmt.Sprintf("%v.%v.%v.%v (%s)", v.Major, v.Minor, v.Patch, v.Build, date)
	ver := fmt.Sprintf("%v.%v.%v (%s)", v.Major, v.Minor, v.Patch, date)
	txt := fmt.Sprintf(tmpl, ver)
	err := ioutil.WriteFile(fversion, []byte(txt), 0644)
	if err != nil {
		return err
	}

	return nil
}
