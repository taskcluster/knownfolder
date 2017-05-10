package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"syscall"
	"unsafe"

	docopt "github.com/docopt/docopt-go"
)

// note, we need to use var not const to avoid compilation failure due to intended overflow
var minusOne = -1

// https://msdn.microsoft.com/en-us/library/windows/desktop/dd378457(v=vs.85).aspx
var (
	knownfolders map[string]*syscall.GUID = map[string]*syscall.GUID{
		"AccountPictures":        {0x008CA0B1, 0x55B4, 0x4C56, [8]byte{0xB8, 0xA8, 0x4D, 0xE4, 0xB2, 0x99, 0xD3, 0xBE}},
		"AddNewPrograms":         {0xDE61D971, 0x5EBC, 0x4F02, [8]byte{0xA3, 0xA9, 0x6C, 0x82, 0x89, 0x5E, 0x5C, 0x04}},
		"AdminTools":             {0x724EF170, 0xA42D, 0x4FEF, [8]byte{0x9F, 0x26, 0xB6, 0x0E, 0x84, 0x6F, 0xBA, 0x4F}},
		"ApplicationShortcuts":   {0xA3918781, 0xE5F2, 0x4890, [8]byte{0xB3, 0xD9, 0xA7, 0xE5, 0x43, 0x32, 0x32, 0x8C}},
		"AppsFolder":             {0x1E87508D, 0x89C2, 0x42F0, [8]byte{0x8A, 0x7E, 0x64, 0x5A, 0x0F, 0x50, 0xCA, 0x58}},
		"AppUpdates":             {0xA305CE99, 0xF527, 0x492B, [8]byte{0x8B, 0x1A, 0x7E, 0x76, 0xFA, 0x98, 0xD6, 0xE4}},
		"CameraRoll":             {0xAB5FB87B, 0x7CE2, 0x4F83, [8]byte{0x91, 0x5D, 0x55, 0x08, 0x46, 0xC9, 0x53, 0x7B}},
		"CDBurning":              {0x9E52AB10, 0xF80D, 0x49DF, [8]byte{0xAC, 0xB8, 0x43, 0x30, 0xF5, 0x68, 0x78, 0x55}},
		"ChangeRemovePrograms":   {0xDF7266AC, 0x9274, 0x4867, [8]byte{0x8D, 0x55, 0x3B, 0xD6, 0x61, 0xDE, 0x87, 0x2D}},
		"CommonAdminTools":       {0xD0384E7D, 0xBAC3, 0x4797, [8]byte{0x8F, 0x14, 0xCB, 0xA2, 0x29, 0xB3, 0x92, 0xB5}},
		"CommonOEMLinks":         {0xC1BAE2D0, 0x10DF, 0x4334, [8]byte{0xBE, 0xDD, 0x7A, 0xA2, 0x0B, 0x22, 0x7A, 0x9D}},
		"CommonPrograms":         {0x0139D44E, 0x6AFE, 0x49F2, [8]byte{0x86, 0x90, 0x3D, 0xAF, 0xCA, 0xE6, 0xFF, 0xB8}},
		"CommonStartMenu":        {0xA4115719, 0xD62E, 0x491D, [8]byte{0xAA, 0x7C, 0xE7, 0x4B, 0x8B, 0xE3, 0xB0, 0x67}},
		"CommonStartup":          {0x82A5EA35, 0xD9CD, 0x47C5, [8]byte{0x96, 0x29, 0xE1, 0x5D, 0x2F, 0x71, 0x4E, 0x6E}},
		"CommonTemplates":        {0xB94237E7, 0x57AC, 0x4347, [8]byte{0x91, 0x51, 0xB0, 0x8C, 0x6C, 0x32, 0xD1, 0xF7}},
		"ComputerFolder":         {0x0AC0837C, 0xBBF8, 0x452A, [8]byte{0x85, 0x0D, 0x79, 0xD0, 0x8E, 0x66, 0x7C, 0xA7}},
		"ConflictFolder":         {0x4BFEFB45, 0x347D, 0x4006, [8]byte{0xA5, 0xBE, 0xAC, 0x0C, 0xB0, 0x56, 0x71, 0x92}},
		"ConnectionsFolder":      {0x6F0CD92B, 0x2E97, 0x45D1, [8]byte{0x88, 0xFF, 0xB0, 0xD1, 0x86, 0xB8, 0xDE, 0xDD}},
		"Contacts":               {0x56784854, 0xC6CB, 0x462B, [8]byte{0x81, 0x69, 0x88, 0xE3, 0x50, 0xAC, 0xB8, 0x82}},
		"ControlPanelFolder":     {0x82A74AEB, 0xAEB4, 0x465C, [8]byte{0xA0, 0x14, 0xD0, 0x97, 0xEE, 0x34, 0x6D, 0x63}},
		"Cookies":                {0x2B0F765D, 0xC0E9, 0x4171, [8]byte{0x90, 0x8E, 0x08, 0xA6, 0x11, 0xB8, 0x4F, 0xF6}},
		"Desktop":                {0xB4BFCC3A, 0xDB2C, 0x424C, [8]byte{0xB0, 0x29, 0x7F, 0xE9, 0x9A, 0x87, 0xC6, 0x41}},
		"DeviceMetadataStore":    {0x5CE4A5E9, 0xE4EB, 0x479D, [8]byte{0xB8, 0x9F, 0x13, 0x0C, 0x02, 0x88, 0x61, 0x55}},
		"Documents":              {0xFDD39AD0, 0x238F, 0x46AF, [8]byte{0xAD, 0xB4, 0x6C, 0x85, 0x48, 0x03, 0x69, 0xC7}},
		"DocumentsLibrary":       {0x7B0DB17D, 0x9CD2, 0x4A93, [8]byte{0x97, 0x33, 0x46, 0xCC, 0x89, 0x02, 0x2E, 0x7C}},
		"Downloads":              {0x374DE290, 0x123F, 0x4565, [8]byte{0x91, 0x64, 0x39, 0xC4, 0x92, 0x5E, 0x46, 0x7B}},
		"Favorites":              {0x1777F761, 0x68AD, 0x4D8A, [8]byte{0x87, 0xBD, 0x30, 0xB7, 0x59, 0xFA, 0x33, 0xDD}},
		"Fonts":                  {0xFD228CB7, 0xAE11, 0x4AE3, [8]byte{0x86, 0x4C, 0x16, 0xF3, 0x91, 0x0A, 0xB8, 0xFE}},
		"Games":                  {0xCAC52C1A, 0xB53D, 0x4EDC, [8]byte{0x92, 0xD7, 0x6B, 0x2E, 0x8A, 0xC1, 0x94, 0x34}},
		"GameTasks":              {0x054FAE61, 0x4DD8, 0x4787, [8]byte{0x80, 0xB6, 0x09, 0x02, 0x20, 0xC4, 0xB7, 0x00}},
		"History":                {0xD9DC8A3B, 0xB784, 0x432E, [8]byte{0xA7, 0x81, 0x5A, 0x11, 0x30, 0xA7, 0x59, 0x63}},
		"HomeGroup":              {0x52528A6B, 0xB9E3, 0x4ADD, [8]byte{0xB6, 0x0D, 0x58, 0x8C, 0x2D, 0xBA, 0x84, 0x2D}},
		"HomeGroupCurrentUser":   {0x9B74B6A3, 0x0DFD, 0x4F11, [8]byte{0x9E, 0x78, 0x5F, 0x78, 0x00, 0xF2, 0xE7, 0x72}},
		"ImplicitAppShortcuts":   {0xBCB5256F, 0x79F6, 0x4CEE, [8]byte{0xB7, 0x25, 0xDC, 0x34, 0xE4, 0x02, 0xFD, 0x46}},
		"InternetCache":          {0x352481E8, 0x33BE, 0x4251, [8]byte{0xBA, 0x85, 0x60, 0x07, 0xCA, 0xED, 0xCF, 0x9D}},
		"InternetFolder":         {0x4D9F7874, 0x4E0C, 0x4904, [8]byte{0x96, 0x7B, 0x40, 0xB0, 0xD2, 0x0C, 0x3E, 0x4B}},
		"Libraries":              {0x1B3EA5DC, 0xB587, 0x4786, [8]byte{0xB4, 0xEF, 0xBD, 0x1D, 0xC3, 0x32, 0xAE, 0xAE}},
		"Links":                  {0xBFB9D5E0, 0xC6A9, 0x404C, [8]byte{0xB2, 0xB2, 0xAE, 0x6D, 0xB6, 0xAF, 0x49, 0x68}},
		"LocalAppData":           {0xF1B32785, 0x6FBA, 0x4FCF, [8]byte{0x9D, 0x55, 0x7B, 0x8E, 0x7F, 0x15, 0x70, 0x91}},
		"LocalAppDataLow":        {0xA520A1A4, 0x1780, 0x4FF6, [8]byte{0xBD, 0x18, 0x16, 0x73, 0x43, 0xC5, 0xAF, 0x16}},
		"LocalizedResourcesDir":  {0x2A00375E, 0x224C, 0x49DE, [8]byte{0xB8, 0xD1, 0x44, 0x0D, 0xF7, 0xEF, 0x3D, 0xDC}},
		"Music":                  {0x4BD8D571, 0x6D19, 0x48D3, [8]byte{0xBE, 0x97, 0x42, 0x22, 0x20, 0x08, 0x0E, 0x43}},
		"MusicLibrary":           {0x2112AB0A, 0xC86A, 0x4FFE, [8]byte{0xA3, 0x68, 0x0D, 0xE9, 0x6E, 0x47, 0x01, 0x2E}},
		"NetHood":                {0xC5ABBF53, 0xE17F, 0x4121, [8]byte{0x89, 0x00, 0x86, 0x62, 0x6F, 0xC2, 0xC9, 0x73}},
		"NetworkFolder":          {0xD20BEEC4, 0x5CA8, 0x4905, [8]byte{0xAE, 0x3B, 0xBF, 0x25, 0x1E, 0xA0, 0x9B, 0x53}},
		"OriginalImages":         {0x2C36C0AA, 0x5812, 0x4B87, [8]byte{0xBF, 0xD0, 0x4C, 0xD0, 0xDF, 0xB1, 0x9B, 0x39}},
		"PhotoAlbums":            {0x69D2CF90, 0xFC33, 0x4FB7, [8]byte{0x9A, 0x0C, 0xEB, 0xB0, 0xF0, 0xFC, 0xB4, 0x3C}},
		"PicturesLibrary":        {0xA990AE9F, 0xA03B, 0x4E80, [8]byte{0x94, 0xBC, 0x99, 0x12, 0xD7, 0x50, 0x41, 0x04}},
		"Pictures":               {0x33E28130, 0x4E1E, 0x4676, [8]byte{0x83, 0x5A, 0x98, 0x39, 0x5C, 0x3B, 0xC3, 0xBB}},
		"Playlists":              {0xDE92C1C7, 0x837F, 0x4F69, [8]byte{0xA3, 0xBB, 0x86, 0xE6, 0x31, 0x20, 0x4A, 0x23}},
		"PrintersFolder":         {0x76FC4E2D, 0xD6AD, 0x4519, [8]byte{0xA6, 0x63, 0x37, 0xBD, 0x56, 0x06, 0x81, 0x85}},
		"PrintHood":              {0x9274BD8D, 0xCFD1, 0x41C3, [8]byte{0xB3, 0x5E, 0xB1, 0x3F, 0x55, 0xA7, 0x58, 0xF4}},
		"Profile":                {0x5E6C858F, 0x0E22, 0x4760, [8]byte{0x9A, 0xFE, 0xEA, 0x33, 0x17, 0xB6, 0x71, 0x73}},
		"ProgramData":            {0x62AB5D82, 0xFDC1, 0x4DC3, [8]byte{0xA9, 0xDD, 0x07, 0x0D, 0x1D, 0x49, 0x5D, 0x97}},
		"ProgramFiles":           {0x905E63B6, 0xC1BF, 0x494E, [8]byte{0xB2, 0x9C, 0x65, 0xB7, 0x32, 0xD3, 0xD2, 0x1A}},
		"ProgramFilesX64":        {0x6D809377, 0x6AF0, 0x444B, [8]byte{0x89, 0x57, 0xA3, 0x77, 0x3F, 0x02, 0x20, 0x0E}},
		"ProgramFilesX86":        {0x7C5A40EF, 0xA0FB, 0x4BFC, [8]byte{0x87, 0x4A, 0xC0, 0xF2, 0xE0, 0xB9, 0xFA, 0x8E}},
		"ProgramFilesCommon":     {0xF7F1ED05, 0x9F6D, 0x47A2, [8]byte{0xAA, 0xAE, 0x29, 0xD3, 0x17, 0xC6, 0xF0, 0x66}},
		"ProgramFilesCommonX64":  {0x6365D5A7, 0x0F0D, 0x45E5, [8]byte{0x87, 0xF6, 0x0D, 0xA5, 0x6B, 0x6A, 0x4F, 0x7D}},
		"ProgramFilesCommonX86":  {0xDE974D24, 0xD9C6, 0x4D3E, [8]byte{0xBF, 0x91, 0xF4, 0x45, 0x51, 0x20, 0xB9, 0x17}},
		"Programs":               {0xA77F5D77, 0x2E2B, 0x44C3, [8]byte{0xA6, 0xA2, 0xAB, 0xA6, 0x01, 0x05, 0x4A, 0x51}},
		"Public":                 {0xDFDF76A2, 0xC82A, 0x4D63, [8]byte{0x90, 0x6A, 0x56, 0x44, 0xAC, 0x45, 0x73, 0x85}},
		"PublicDesktop":          {0xC4AA340D, 0xF20F, 0x4863, [8]byte{0xAF, 0xEF, 0xF8, 0x7E, 0xF2, 0xE6, 0xBA, 0x25}},
		"PublicDocuments":        {0xED4824AF, 0xDCE4, 0x45A8, [8]byte{0x81, 0xE2, 0xFC, 0x79, 0x65, 0x08, 0x36, 0x34}},
		"PublicDownloads":        {0x3D644C9B, 0x1FB8, 0x4F30, [8]byte{0x9B, 0x45, 0xF6, 0x70, 0x23, 0x5F, 0x79, 0xC0}},
		"PublicGameTasks":        {0xDEBF2536, 0xE1A8, 0x4C59, [8]byte{0xB6, 0xA2, 0x41, 0x45, 0x86, 0x47, 0x6A, 0xEA}},
		"PublicLibraries":        {0x48DAF80B, 0xE6CF, 0x4F4E, [8]byte{0xB8, 0x00, 0x0E, 0x69, 0xD8, 0x4E, 0xE3, 0x84}},
		"PublicMusic":            {0x3214FAB5, 0x9757, 0x4298, [8]byte{0xBB, 0x61, 0x92, 0xA9, 0xDE, 0xAA, 0x44, 0xFF}},
		"PublicPictures":         {0xB6EBFB86, 0x6907, 0x413C, [8]byte{0x9A, 0xF7, 0x4F, 0xC2, 0xAB, 0xF0, 0x7C, 0xC5}},
		"PublicRingtones":        {0xE555AB60, 0x153B, 0x4D17, [8]byte{0x9F, 0x04, 0xA5, 0xFE, 0x99, 0xFC, 0x15, 0xEC}},
		"PublicUserTiles":        {0x0482AF6C, 0x08F1, 0x4C34, [8]byte{0x8C, 0x90, 0xE1, 0x7E, 0xC9, 0x8B, 0x1E, 0x17}},
		"PublicVideos":           {0x2400183A, 0x6185, 0x49FB, [8]byte{0xA2, 0xD8, 0x4A, 0x39, 0x2A, 0x60, 0x2B, 0xA3}},
		"QuickLaunch":            {0x52A4F021, 0x7B75, 0x48A9, [8]byte{0x9F, 0x6B, 0x4B, 0x87, 0xA2, 0x10, 0xBC, 0x8F}},
		"Recent":                 {0xAE50C081, 0xEBD2, 0x438A, [8]byte{0x86, 0x55, 0x8A, 0x09, 0x2E, 0x34, 0x98, 0x7A}},
		"RecordedTVLibrary":      {0x1A6FDBA2, 0xF42D, 0x4358, [8]byte{0xA7, 0x98, 0xB7, 0x4D, 0x74, 0x59, 0x26, 0xC5}},
		"RecycleBinFolder":       {0xB7534046, 0x3ECB, 0x4C18, [8]byte{0xBE, 0x4E, 0x64, 0xCD, 0x4C, 0xB7, 0xD6, 0xAC}},
		"ResourceDir":            {0x8AD10C31, 0x2ADB, 0x4296, [8]byte{0xA8, 0xF7, 0xE4, 0x70, 0x12, 0x32, 0xC9, 0x72}},
		"Ringtones":              {0xC870044B, 0xF49E, 0x4126, [8]byte{0xA9, 0xC3, 0xB5, 0x2A, 0x1F, 0xF4, 0x11, 0xE8}},
		"RoamingAppData":         {0x3EB685DB, 0x65F9, 0x4CF6, [8]byte{0xA0, 0x3A, 0xE3, 0xEF, 0x65, 0x72, 0x9F, 0x3D}},
		"RoamedTileImages":       {0xAAA8D5A5, 0xF1D6, 0x4259, [8]byte{0xBA, 0xA8, 0x78, 0xE7, 0xEF, 0x60, 0x83, 0x5E}},
		"RoamingTiles":           {0x00BCFC5A, 0xED94, 0x4E48, [8]byte{0x96, 0xA1, 0x3F, 0x62, 0x17, 0xF2, 0x19, 0x90}},
		"SampleMusic":            {0xB250C668, 0xF57D, 0x4EE1, [8]byte{0xA6, 0x3C, 0x29, 0x0E, 0xE7, 0xD1, 0xAA, 0x1F}},
		"SamplePictures":         {0xC4900540, 0x2379, 0x4C75, [8]byte{0x84, 0x4B, 0x64, 0xE6, 0xFA, 0xF8, 0x71, 0x6B}},
		"SamplePlaylists":        {0x15CA69B3, 0x30EE, 0x49C1, [8]byte{0xAC, 0xE1, 0x6B, 0x5E, 0xC3, 0x72, 0xAF, 0xB5}},
		"SampleVideos":           {0x859EAD94, 0x2E85, 0x48AD, [8]byte{0xA7, 0x1A, 0x09, 0x69, 0xCB, 0x56, 0xA6, 0xCD}},
		"SavedGames":             {0x4C5C32FF, 0xBB9D, 0x43B0, [8]byte{0xB5, 0xB4, 0x2D, 0x72, 0xE5, 0x4E, 0xAA, 0xA4}},
		"SavedPictures":          {0x3B193882, 0xD3AD, 0x4EAB, [8]byte{0x96, 0x5A, 0x69, 0x82, 0x9D, 0x1F, 0xB5, 0x9F}},
		"SavedPicturesLibrary":   {0xE25B5812, 0xBE88, 0x4BD9, [8]byte{0x94, 0xB0, 0x29, 0x23, 0x34, 0x77, 0xB6, 0xC3}},
		"SavedSearches":          {0x7D1D3A04, 0xDEBB, 0x4115, [8]byte{0x95, 0xCF, 0x2F, 0x29, 0xDA, 0x29, 0x20, 0xDA}},
		"Screenshots":            {0xB7BEDE81, 0xDF94, 0x4682, [8]byte{0xA7, 0xD8, 0x57, 0xA5, 0x26, 0x20, 0xB8, 0x6F}},
		"SEARCH_CSC":             {0xEE32E446, 0x31CA, 0x4ABA, [8]byte{0x81, 0x4F, 0xA5, 0xEB, 0xD2, 0xFD, 0x6D, 0x5E}},
		"SearchHistory":          {0x0D4C3DB6, 0x03A3, 0x462F, [8]byte{0xA0, 0xE6, 0x08, 0x92, 0x4C, 0x41, 0xB5, 0xD4}},
		"SearchHome":             {0x190337D1, 0xB8CA, 0x4121, [8]byte{0xA6, 0x39, 0x6D, 0x47, 0x2D, 0x16, 0x97, 0x2A}},
		"SEARCH_MAPI":            {0x98EC0E18, 0x2098, 0x4D44, [8]byte{0x86, 0x44, 0x66, 0x97, 0x93, 0x15, 0xA2, 0x81}},
		"SearchTemplates":        {0x7E636BFE, 0xDFA9, 0x4D5E, [8]byte{0xB4, 0x56, 0xD7, 0xB3, 0x98, 0x51, 0xD8, 0xA9}},
		"SendTo":                 {0x8983036C, 0x27C0, 0x404B, [8]byte{0x8F, 0x08, 0x10, 0x2D, 0x10, 0xDC, 0xFD, 0x74}},
		"SidebarDefaultParts":    {0x7B396E54, 0x9EC5, 0x4300, [8]byte{0xBE, 0x0A, 0x24, 0x82, 0xEB, 0xAE, 0x1A, 0x26}},
		"SidebarParts":           {0xA75D362E, 0x50FC, 0x4FB7, [8]byte{0xAC, 0x2C, 0xA8, 0xBE, 0xAA, 0x31, 0x44, 0x93}},
		"SkyDrive":               {0xA52BBA46, 0xE9E1, 0x435F, [8]byte{0xB3, 0xD9, 0x28, 0xDA, 0xA6, 0x48, 0xC0, 0xF6}},
		"SkyDriveCameraRoll":     {0x767E6811, 0x49CB, 0x4273, [8]byte{0x87, 0xC2, 0x20, 0xF3, 0x55, 0xE1, 0x08, 0x5B}},
		"SkyDriveDocuments":      {0x24D89E24, 0x2F19, 0x4534, [8]byte{0x9D, 0xDE, 0x6A, 0x66, 0x71, 0xFB, 0xB8, 0xFE}},
		"SkyDrivePictures":       {0x339719B5, 0x8C47, 0x4894, [8]byte{0x94, 0xC2, 0xD8, 0xF7, 0x7A, 0xDD, 0x44, 0xA6}},
		"StartMenu":              {0x625B53C3, 0xAB48, 0x4EC1, [8]byte{0xBA, 0x1F, 0xA1, 0xEF, 0x41, 0x46, 0xFC, 0x19}},
		"Startup":                {0xB97D20BB, 0xF46A, 0x4C97, [8]byte{0xBA, 0x10, 0x5E, 0x36, 0x08, 0x43, 0x08, 0x54}},
		"SyncManagerFolder":      {0x43668BF8, 0xC14E, 0x49B2, [8]byte{0x97, 0xC9, 0x74, 0x77, 0x84, 0xD7, 0x84, 0xB7}},
		"SyncResultsFolder":      {0x289A9A43, 0xBE44, 0x4057, [8]byte{0xA4, 0x1B, 0x58, 0x7A, 0x76, 0xD7, 0xE7, 0xF9}},
		"SyncSetupFolder":        {0x0F214138, 0xB1D3, 0x4A90, [8]byte{0xBB, 0xA9, 0x27, 0xCB, 0xC0, 0xC5, 0x38, 0x9A}},
		"System":                 {0x1AC14E77, 0x02E7, 0x4E5D, [8]byte{0xB7, 0x44, 0x2E, 0xB1, 0xAE, 0x51, 0x98, 0xB7}},
		"SystemX86":              {0xD65231B0, 0xB2F1, 0x4857, [8]byte{0xA4, 0xCE, 0xA8, 0xE7, 0xC6, 0xEA, 0x7D, 0x27}},
		"Templates":              {0xA63293E8, 0x664E, 0x48DB, [8]byte{0xA0, 0x79, 0xDF, 0x75, 0x9E, 0x05, 0x09, 0xF7}},
		"UserPinned":             {0x9E3995AB, 0x1F9C, 0x4F13, [8]byte{0xB8, 0x27, 0x48, 0xB2, 0x4B, 0x6C, 0x71, 0x74}},
		"UserProfiles":           {0x0762D272, 0xC50A, 0x4BB0, [8]byte{0xA3, 0x82, 0x69, 0x7D, 0xCD, 0x72, 0x9B, 0x80}},
		"UserProgramFiles":       {0x5CD7AEE2, 0x2219, 0x4A67, [8]byte{0xB8, 0x5D, 0x6C, 0x9C, 0xE1, 0x56, 0x60, 0xCB}},
		"UserProgramFilesCommon": {0xBCBD3057, 0xCA5C, 0x4622, [8]byte{0xB4, 0x2D, 0xBC, 0x56, 0xDB, 0x0A, 0xE5, 0x16}},
		"UsersFiles":             {0xF3CE0F7C, 0x4901, 0x4ACC, [8]byte{0x86, 0x48, 0xD5, 0xD4, 0x4B, 0x04, 0xEF, 0x8F}},
		"UsersLibraries":         {0xA302545D, 0xDEFF, 0x464B, [8]byte{0xAB, 0xE8, 0x61, 0xC8, 0x64, 0x8D, 0x93, 0x9B}},
		"Videos":                 {0x18989B1D, 0x99B5, 0x455B, [8]byte{0x84, 0x1C, 0xAB, 0x7C, 0x74, 0xE4, 0xDD, 0xFC}},
		"VideosLibrary":          {0x491E922F, 0x5643, 0x4AF4, [8]byte{0xA7, 0xEB, 0x4E, 0x7A, 0x13, 0x8D, 0x81, 0x74}},
		"Windows":                {0xF38BF404, 0x1D43, 0x42F2, [8]byte{0x93, 0x05, 0x67, 0xDE, 0x0B, 0x28, 0xFC, 0x23}},
	}
)

type ProfileInfo struct {
	Size        uint32
	Flags       uint32
	Username    *uint16
	ProfilePath *uint16
	DefaultPath *uint16
	ServerName  *uint16
	PolicyPath  *uint16
	Profile     syscall.Handle
}

const (
	PI_NOUI = 1

	LOGON32_PROVIDER_DEFAULT = 0
	LOGON32_PROVIDER_WINNT35 = 1
	LOGON32_PROVIDER_WINNT40 = 2
	LOGON32_PROVIDER_WINNT50 = 3

	LOGON32_LOGON_INTERACTIVE       = 2
	LOGON32_LOGON_NETWORK           = 3
	LOGON32_LOGON_BATCH             = 4
	LOGON32_LOGON_SERVICE           = 5
	LOGON32_LOGON_UNLOCK            = 7
	LOGON32_LOGON_NETWORK_CLEARTEXT = 8
	LOGON32_LOGON_NEW_CREDENTIALS   = 9
)

var (
	version = "1.0.0"
	usage   = `
knownfolder

knownfolder allows you to get and set known folder locations on Windows.

See https://msdn.microsoft.com/en-us/library/windows/desktop/dd378457(v=vs.85).aspx

  Usage:
    knownfolder set [-d|-u USERNAME -p PASSWORD] FOLDER LOCATION
    knownfolder get [-d|-u USERNAME -p PASSWORD] FOLDER
    knownfolder list
    knownfolder -h|--help
    knownfolder --version

  Targets:
    set          Set a folder location. You need to run this command as the user concerned, for
                 USER based settings, otherwise -d will apply the setting for the default user.
    get          Retrieve a folder location. You need to run this command as the user concerned,
                 for USER based settings, otherwise -d will apply the setting for the default user.
    list         List all possible values for FOLDER.

  Options:
    -d           Set/get known folder for the default user profile, rather than an existing user.
    FOLDER       The folder name, as per the Constants shown in
                 https://msdn.microsoft.com/en-us/library/windows/desktop/dd378457(v=vs.85).aspx
    LOCATION     The full file system path to set the given FOLDER location to.
    USERNAME     The username of the user you wish to set/get the known folder for, if different
                 to the user running the knownfolder command.
    PASSWORD     The password of the user you wish to set/get the known folder for, if different
                 to the user running the knownfolder command.

  Examples:

    C:\> knownfolder set RoamingAppData "D:\Users\Pete\AppData\Roaming"
    C:\> knownfolder list
    C:\> knownfolder get LocalAppData
    C:\> knownfolder --help
    C:\> knownfolder --version
`
	advapi32                 = syscall.NewLazyDLL("advapi32.dll")
	shell32                  = syscall.NewLazyDLL("shell32.dll")
	ole32                    = syscall.NewLazyDLL("ole32.dll")
	userenv                  = syscall.NewLazyDLL("userenv.dll")
	procSHGetKnownFolderPath = shell32.NewProc("SHGetKnownFolderPath")
	procCoTaskMemFree        = ole32.NewProc("CoTaskMemFree")
	procSHSetKnownFolderPath = shell32.NewProc("SHSetKnownFolderPath")
	procLoadUserProfileW     = userenv.NewProc("LoadUserProfileW")
	procLogonUserW           = advapi32.NewProc("LogonUserW")
	procUnloadUserProfile    = userenv.NewProc("UnloadUserProfile")
)

func main() {
	arguments, err := docopt.Parse(usage, nil, true, "knownfolders "+version, false, true)
	if err != nil {
		log.Fatalf("Error parsing command line arguments!")
	}

	switch {
	case arguments["set"]:
		location := arguments["LOCATION"].(string)
		folder := arguments["FOLDER"].(string)
		if knownfolders[folder] == nil {
			log.Fatalf(`Unknown folder "%v"`, folder)
		}
		hUser := syscall.Handle(0)
		if arguments["-d"].(bool) {
			// intentionally overflow minusOne to uintptr 0xFFFF.... here
			hUser = syscall.Handle(minusOne)
		} else if arguments["USERNAME"].(bool) {
			var profileInfo *ProfileInfo
			hUser, profileInfo = InteractiveLogonUser(arguments["USERNAME"].(string), arguments["PASSWORD"].(string))
			defer LogoffUser(hUser, profileInfo)
		}
		err := SetFolder(hUser, knownfolders[folder], location)
		if err != nil {
			log.Fatalf("Could not set folder location %v=%v\n%v", folder, location, err)
		}
		fmt.Printf("%v=%v", folder, location)
	case arguments["get"]:
		folder := arguments["FOLDER"].(string)
		if knownfolders[folder] == nil {
			log.Fatalf(`Unknown folder "%v"`, folder)
		}
		hUser := syscall.Handle(0)
		if arguments["-d"].(bool) {
			// intentionally overflow minusOne to uintptr 0xFFFF.... here
			hUser = syscall.Handle(minusOne)
		} else if arguments["-u"].(bool) {
			var profileInfo *ProfileInfo
			hUser, profileInfo = InteractiveLogonUser(arguments["USERNAME"].(string), arguments["PASSWORD"].(string))
			defer LogoffUser(hUser, profileInfo)
		}
		value, err := GetFolder(hUser, knownfolders[folder])
		if err != nil {
			log.Fatalf("Could not retrieve folder %v:\n%v", folder, err)
		}
		fmt.Println(value)
	case arguments["list"]:
		err := ListFolders()
		if err != nil {
			log.Fatalf("Could not list folders:\n%v", err)
		}
	}
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/bb762188(v=vs.85).aspx
func SHGetKnownFolderPath(rfid *syscall.GUID, dwFlags uint32, hToken syscall.Handle, pszPath *uintptr) (err error) {
	r0, _, _ := procSHGetKnownFolderPath.Call(
		uintptr(unsafe.Pointer(rfid)),
		uintptr(dwFlags),
		uintptr(hToken),
		uintptr(unsafe.Pointer(pszPath)),
	)
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/bb762249(v=vs.85).aspx
func SHSetKnownFolderPath(
	rfid *syscall.GUID, // REFKNOWNFOLDERID
	dwFlags uint32, // DWORD
	hToken syscall.Handle, // HANDLE
	pszPath *uint16, // PCWSTR
) (err error) {
	r1, _, _ := procSHSetKnownFolderPath.Call(
		uintptr(unsafe.Pointer(rfid)),
		uintptr(dwFlags),
		uintptr(hToken),
		uintptr(unsafe.Pointer(pszPath)),
	)
	if r1 != 0 {
		err = syscall.Errno(r1)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms680722(v=vs.85).aspx
// Note: the system call returns no value, so we can't check for an error
func CoTaskMemFree(pv uintptr) {
	procCoTaskMemFree.Call(uintptr(pv))
}

func GetFolder(hUser syscall.Handle, folder *syscall.GUID) (value string, err error) {
	var path uintptr
	err = SHGetKnownFolderPath(folder, 0, hUser, &path)

	if err != nil {
		return
	}
	// CoTaskMemFree system call has no return value, so can't check for error
	defer CoTaskMemFree(path)
	value = syscall.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(path))[:])
	return
}

func SetFolder(hUser syscall.Handle, folder *syscall.GUID, value string) (err error) {
	var s *uint16
	s, err = syscall.UTF16PtrFromString(value)
	if err != nil {
		return
	}
	return SHSetKnownFolderPath(folder, 0, hUser, s)
}

func ListFolders() (err error) {
	keys := make([]string, 0, len(knownfolders))
	for key := range knownfolders {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Println(key)
	}
	return
}

func InteractiveLogonUser(username, password string) (user syscall.Handle, pinfo *ProfileInfo) {

	name, err := syscall.UTF16PtrFromString(username)
	if err != nil {
		panic(err)
	}

	pinfo = &ProfileInfo{
		Size:     uint32(unsafe.Sizeof(*pinfo)),
		Flags:    PI_NOUI,
		Username: name,
	}

	// first log on user ....

	user, err = LogonUser(
		syscall.StringToUTF16Ptr(username),
		syscall.StringToUTF16Ptr("."),
		syscall.StringToUTF16Ptr(password),
		LOGON32_LOGON_INTERACTIVE,
		LOGON32_PROVIDER_DEFAULT,
	)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// now load user profile ....

	err = LoadUserProfile(user, pinfo)
	if err != nil {
		panic(err)
	}
	return
}

func LogonUser(username *uint16, domain *uint16, password *uint16, logonType uint32, logonProvider uint32) (token syscall.Handle, err error) {
	r1, _, e1 := procLogonUserW.Call(
		uintptr(unsafe.Pointer(username)),
		uintptr(unsafe.Pointer(domain)),
		uintptr(unsafe.Pointer(password)),
		uintptr(logonType),
		uintptr(logonProvider),
		uintptr(unsafe.Pointer(&token)))
	runtime.KeepAlive(username)
	runtime.KeepAlive(domain)
	runtime.KeepAlive(password)
	if int(r1) == 0 {
		return syscall.InvalidHandle, os.NewSyscallError("LogonUser", e1)
	}
	return
}

func LogoffUser(user syscall.Handle, pinfo *ProfileInfo) {
	defer syscall.Close(user)
	defer func() {
		if pinfo.Profile != syscall.Handle(0) && pinfo.Profile != syscall.InvalidHandle {
			for {
				err := UnloadUserProfile(user, pinfo.Profile)
				if err == nil {
					break
				}
				log.Printf("%v", err)
			}
		}
	}()
}

func LoadUserProfile(token syscall.Handle, pinfo *ProfileInfo) error {
	r1, _, e1 := procLoadUserProfileW.Call(
		uintptr(token),
		uintptr(unsafe.Pointer(pinfo)))
	runtime.KeepAlive(pinfo)
	if int(r1) == 0 {
		return os.NewSyscallError("LoadUserProfile", e1)
	}
	return nil
}

func UnloadUserProfile(token, profile syscall.Handle) error {
	if r1, _, e1 := procUnloadUserProfile.Call(
		uintptr(token),
		uintptr(profile)); int(r1) == 0 {
		return os.NewSyscallError("UnloadUserProfile", e1)
	}
	return nil
}
