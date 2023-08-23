# patchprep
Prepare a patch for eq patches

## Usage

Note it is important that your obtained copy of EQ has original modified date stamps. If they aren't, you may need to download a fresh copy to compare against. You can verify this by going into your EQ directory and looking at the modified dates, most are 2009 or so.

- Download patchprep.exe from [releases](https://github.com/xackery/patchprep/releases)
- Copy it to your EQ directory
- Double click patchprep.exe
- It will scan all your EQ files, and look for any files modified within the last 2 years and aren't filtered in patchprep.txt.
- Review the files it copies to the patch subfolder
- Modify what files to exclude in the patchprep.txt file using .github like patterns
- Run it again
- Review what files are generated, when they look good, copy them to [eqemupatcher](https://github.com/xackery/eqemupatcher/) or [launcheq](https://github.com/xackery/launcheq/)'s rof folder
- Push up your changes
- Profit
