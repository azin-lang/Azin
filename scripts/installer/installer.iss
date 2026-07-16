#define AzAppVersion "0.2.1"

[Setup]
AppName=Azin
AppVersion={#AzAppVersion}
AppPublisher=Azin maintainers
AppPublisherURL=https://github.com/azin-lang/Azin
AppSupportURL=https://github.com/azin-lang/Azin/issues

DefaultDirName={autopf}\Azin
DefaultGroupName=Azin
DisableProgramGroupPage=yes

OutputDir=..\installer
OutputBaseFilename=Azin-setup-{#AzAppVersion}

Compression=lzma2
SolidCompression=yes

PrivilegesRequired=lowest
PrivilegesRequiredOverridesAllowed=dialog
ArchitecturesInstallIn64BitMode=x64compatible

ChangesEnvironment=yes

SetupIconFile=azin.ico
UninstallDisplayIcon={app}\azin.ico


[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"


[Files]
Source: "..\..\build\azc.exe"; DestDir: "{app}"; DestName: "azc.exe"; Flags: ignoreversion
Source: "azin.ico"; DestDir: "{app}"; Flags: ignoreversion


[Icons]
Name: "{group}\Azin"; Filename: "{app}\azc.exe"; IconFilename: "{app}\azin.ico"
Name: "{group}\Uninstall Azin"; Filename: "{uninstallexe}"


;; https://github.com/brysonak/buf/blob/main/Install/install.iss#L42
[Code]

const
  EnvironmentKey = 'Environment';
  WM_SETTINGCHANGE = $001A;
  SMTO_ABORTIFHUNG = $0002;

function SendMessageTimeout(
  hWnd: LongInt;
  Msg: LongWord;
  wParam: LongInt;
  lParam: LongInt;
  fuFlags: LongWord;
  uTimeout: LongWord;
  var lpdwResult: LongWord
): LongWord;
  external 'SendMessageTimeoutW@user32.dll stdcall';

procedure RefreshEnvironment;
var
  MsgResult: LongWord;
begin
  SendMessageTimeout(
    HWND_BROADCAST,
    WM_SETTINGCHANGE,
    0,
    0,
    SMTO_ABORTIFHUNG,
    5000,
    MsgResult
  );
end;

function PathContains(Path, Dir: string): Boolean;
begin
  Result := Pos(';' + Lowercase(Dir) + ';',
    ';' + Lowercase(Path) + ';') > 0;
end;

procedure AddToPath(Dir: string);
var
  Path: string;
begin
  if not RegQueryStringValue(HKCU, EnvironmentKey, 'Path', Path) then
    Path := '';

  if not PathContains(Path, Dir) then
  begin
    if (Path <> '') and (Path[Length(Path)] <> ';') then
      Path := Path + ';';

    Path := Path + Dir;

    RegWriteExpandStringValue(HKCU, EnvironmentKey, 'Path', Path);
  end;
end;

procedure RemoveFromPath(Dir: string);
var
  Path: string;
begin
  if not RegQueryStringValue(HKCU, EnvironmentKey, 'Path', Path) then
    Exit;

  StringChangeEx(Path, ';' + Dir, '', True);
  StringChangeEx(Path, Dir + ';', '', True);
  StringChangeEx(Path, Dir, '', True);

  while Pos(';;', Path) > 0 do
    StringChangeEx(Path, ';;', ';', True);

  if (Length(Path) > 0) and (Path[1] = ';') then
    Delete(Path, 1, 1);

  if (Length(Path) > 0) and (Path[Length(Path)] = ';') then
    Delete(Path, Length(Path), 1);

  RegWriteExpandStringValue(HKCU, EnvironmentKey, 'Path', Path);
end;

procedure CurStepChanged(CurStep: TSetupStep);
begin
  if CurStep = ssPostInstall then
  begin
    AddToPath(ExpandConstant('{app}'));
    RefreshEnvironment;
  end;
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usPostUninstall then
  begin
    RemoveFromPath(ExpandConstant('{app}'));
    RefreshEnvironment;
  end;
end;