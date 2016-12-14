package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/taskcluster/runlib/win32"
)

var (
	version = "1.0.0"
	usage   = `
knownfolder

knownfolder allows you to get and set known folder locations on Windows.

See https://msdn.microsoft.com/en-us/library/windows/desktop/dd378457(v=vs.85).aspx

  Usage:
    knownfolder set FOLDER LOCATION
    knownfolder get FOLDER
    knownfolder list
    knownfolder -h|--help
    knownfolder --version

  Targets:
    set          Set a folder location. You need to run this command as the user concerned, for
                 USER based settings.
    get          Retrieve a folder location. You need to run this command as the user concerned,
                 for USER based settings.
    list         List all possible values for FOLDER.

  Options:
    FOLDER       The folder name, as per the Constants shown in
                 https://msdn.microsoft.com/en-us/library/windows/desktop/dd378457(v=vs.85).aspx
    LOCATION     The full file system path to set the given FOLDER location to.

  Examples:

    C:\> knownfolder set FOLDERID_RoamingAppData "D:\Users\Pete\AppData\Roaming"
    C:\> knownfolder list
    C:\> knownfolder get FOLDERID_LocalAppData
    C:\> knownfolder --help
    C:\> knownfolder --version
`
	modShell32               = syscall.NewLazyDLL("Shell32.dll")
	modOle32                 = syscall.NewLazyDLL("Ole32.dll")
	procSHGetKnownFolderPath = modShell32.NewProc("SHGetKnownFolderPath")
	procCoTaskMemFree        = modOle32.NewProc("CoTaskMemFree")
)

func SHGetKnownFolderPath(rfid *syscall.GUID, dwFlags uint32, hToken syscall.Handle, pszPath *uintptr) (retval error) {
	r0, _, _ := syscall.Syscall6(procSHGetKnownFolderPath.Addr(), 4, uintptr(unsafe.Pointer(rfid)), uintptr(dwFlags), uintptr(hToken), uintptr(unsafe.Pointer(pszPath)), 0, 0)
	if r0 != 0 {
		retval = syscall.Errno(r0)
	}
	return
}

func CoTaskMemFree(pv uintptr) {
	syscall.Syscall(procCoTaskMemFree.Addr(), 1, uintptr(pv), 0, 0)
	return
}

func GetFolder(folder *syscall.GUID) (string, error) {
	var path uintptr
	err := SHGetKnownFolderPath(folder, 0, 0, &path)
	if err != nil {
		return "", err
	}
	defer CoTaskMemFree(path)
	value := syscall.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(path))[:])
	return value, nil
}

func SetFolder(folder *syscall.GUID, value string) (err error) {
	var s *uint16
	s, err = syscall.UTF16PtrFromString(value)
	if err != nil {
		return err
	}
	return win32.SHSetKnownFolderPath(folder, 0, 0, s)
}

func main() {
	folder, err := GetFolder(&win32.FOLDERID_Profile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Original folder:", folder)
	err = SetFolder(&win32.FOLDERID_Profile, "C:\\Users\\%USERNAME%")
	if err != nil {
		fmt.Println(err)
		return
	}
	folder, err = GetFolder(&win32.FOLDERID_Profile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Modified folder:", folder)
}
