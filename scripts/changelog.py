#!/usr/bin/env python3

import re
import sys
import os.path
from pathlib import Path
from typing import List, Optional

from git_changelog import templates
from git_changelog.build import Changelog
from git_changelog.providers import GitHub, RefDef, RefRe

JIRA_URL = "https://issues.redhat.com"
JIRA_PROJECT = "HMSPROV"

GitHub.REF["issues"] = RefDef(
                           regex=re.compile(RefRe.BB + RefRe.ID.format(symbol=JIRA_PROJECT + "-"), re.I),
                           url_string=(JIRA_URL + "/browse/" + JIRA_PROJECT + "-{ref}"),
                       )

def main(args: Optional[List[str]] = None) -> int:
  provider = GitHub("RHEnVision", "provisioning-backend", url="https://github.com")
  path = Path(os.path.abspath(__file__))
  repo = path.parent.parent
  changelog = Changelog(repo, provider=provider, style="angular")
  template = templates.get_template("angular")
  rendered = template.render(changelog=changelog)

  # sys.stdout.write(rendered)
  with open(repo.joinpath("CHANGELOG.md"), "w") as stream:
    stream.write(rendered)

  return 0

if __name__ == '__main__':
  sys.exit(main(sys.argv[1:]))
