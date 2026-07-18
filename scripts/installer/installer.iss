#define AzAppVersion "0.2.1"

[Setup]
AppId={{CA1B358E-4F89-412E-B278-72C2F9B983BD}
AppName=Azin
AppVersion={#AzAppVersion}
AppPublisher=Azin Project
AppPublisherURL=https://github.com/azin-lang/Azin
AppSupportURL=https://github.com/azin-lang/Azin/issues
AppUpdatesURL=https://github.com/azin-lang/Azin/releases

DefaultDirName={localappdata}\Programs\Azin
DefaultGroupName=Azin
DisableProgramGroupPage=yes

LicenseFile=..\..\LICENSE

VersionInfoVersion={#AzAppVersion}
VersionInfoCompany=Azin Project
VersionInfoDescription=Azin compiler installer
VersionInfoCopyright=Azin Project

OutputDir=..\installer
OutputBaseFilename=Azin-setup-{#AzAppVersion}

Compression=lzma2
SolidCompression=yes

PrivilegesRequired=lowest
PrivilegesRequiredOverridesAllowed=dialog
ArchitecturesInstallIn64BitMode=x64compatible

ChangesEnvironment=yes

WizardStyle=modern


[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"


[Files]
Source: "..\..\build\azc.exe"; DestDir: "{app}"; DestName: "azc.exe"; Flags: ignoreversion


[Icons]
Name: "{group}\Azin"; Filename: "{app}\azc.exe"
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

procedure RemoveFromPath(const Dir: string);
var
  Path, NewPath, Entry: string;
  PosSep: Integer;
begin
  if not RegQueryStringValue(HKCU, EnvironmentKey, 'Path', Path) then
    Exit;

  NewPath := '';

  while Path <> '' do
  begin
    PosSep := Pos(';', Path);

    if PosSep = 0 then
    begin
      Entry := Path;
      Path := '';
    end
    else
    begin
      Entry := Copy(Path, 1, PosSep - 1);
      Delete(Path, 1, PosSep);
    end;

    if CompareText(Entry, Dir) <> 0 then
    begin
      if NewPath <> '' then
        NewPath := NewPath + ';';

      NewPath := NewPath + Entry;
    end;
  end;

  RegWriteExpandStringValue(HKCU, EnvironmentKey, 'Path', NewPath);
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