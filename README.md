# knownfolder

This is a command line utility for getting/setting known folders on windows.

See https://msdn.microsoft.com/en-us/library/windows/desktop/dd378457(v=vs.85).aspx for more info.

## Example usage

### Getting help

```
C:\>knownfolder --help
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

    C:\> knownfolder set RoamingAppData "D:\Users\Pete\AppData\Roaming"
    C:\> knownfolder list
    C:\> knownfolder get LocalAppData
    C:\> knownfolder --help
    C:\> knownfolder --version
```

### Setting folder location

```
C:\>knownfolder set RoamingAppData "D:\Users\Pete\AppData\Roaming"
RoamingAppData=D:\Users\Pete\AppData\Roaming
```

### Retrieving folder location

```
C:\>knownfolder get LocalAppData
C:\Tasks\Stuff\AppData\Local

```

### Querying version of knownfolder

```
C:\>knownfolder --version
knownfolders 1.0.0
```

### Listing known folders

```
C:\>knownfolder list
AccountPictures
AddNewPrograms
AdminTools
AppUpdates
ApplicationShortcuts
AppsFolder
CDBurning
CameraRoll
ChangeRemovePrograms
CommonAdminTools
CommonOEMLinks
CommonPrograms
CommonStartMenu
CommonStartup
CommonTemplates
ComputerFolder
ConflictFolder
ConnectionsFolder
Contacts
ControlPanelFolder
Cookies
Desktop
DeviceMetadataStore
Documents
DocumentsLibrary
Downloads
Favorites
Fonts
GameTasks
Games
History
HomeGroup
HomeGroupCurrentUser
ImplicitAppShortcuts
InternetCache
InternetFolder
Libraries
Links
LocalAppData
LocalAppDataLow
LocalizedResourcesDir
Music
MusicLibrary
NetHood
NetworkFolder
OriginalImages
PhotoAlbums
Pictures
PicturesLibrary
Playlists
PrintHood
PrintersFolder
Profile
ProgramData
ProgramFiles
ProgramFilesCommon
ProgramFilesCommonX64
ProgramFilesCommonX86
ProgramFilesX64
ProgramFilesX86
Programs
Public
PublicDesktop
PublicDocuments
PublicDownloads
PublicGameTasks
PublicLibraries
PublicMusic
PublicPictures
PublicRingtones
PublicUserTiles
PublicVideos
QuickLaunch
Recent
RecordedTVLibrary
RecycleBinFolder
ResourceDir
Ringtones
RoamedTileImages
RoamingAppData
RoamingTiles
SEARCH_CSC
SEARCH_MAPI
SampleMusic
SamplePictures
SamplePlaylists
SampleVideos
SavedGames
SavedPictures
SavedPicturesLibrary
SavedSearches
Screenshots
SearchHistory
SearchHome
SearchTemplates
SendTo
SidebarDefaultParts
SidebarParts
SkyDrive
SkyDriveCameraRoll
SkyDriveDocuments
SkyDrivePictures
StartMenu
Startup
SyncManagerFolder
SyncResultsFolder
SyncSetupFolder
System
SystemX86
Templates
UserPinned
UserProfiles
UserProgramFiles
UserProgramFilesCommon
UsersFiles
UsersLibraries
Videos
VideosLibrary
Windows
```
