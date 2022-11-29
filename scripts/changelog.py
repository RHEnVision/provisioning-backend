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
JIRA_ISSUE_SYMBOL = "HMSPROV-"

GitHub.REF["issues"] = RefDef(
                           regex=re.compile(RefRe.BB + RefRe.ID.format(symbol=JIRA_ISSUE_SYMBOL), re.I),
                           url_string=(JIRA_URL + "/browse/" + JIRA_PROJECT + "-{ref}"),
                       )

class SingleCommitLog(Changelog):
  def get_log(self) -> str:
    """Get the `git log` output limited to single commit.

    Returns:
        The output of the `git log` command, with a particular format.
    """
    return self.run_git("log", "-1", "--date=unix", "--format=" + self.FORMAT)

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

def check_commit(length_limit = 70) -> int:
  provider = GitHub("RHEnVision", "provisioning-backend", url="https://github.com")
  path = Path(os.path.abspath(__file__))
  repo = path.parent.parent
  changelog = SingleCommitLog(repo, provider=provider, style="angular")

  commit = changelog.commits[0]
  if len(commit.subject) > length_limit:
    sys.stderr.write("ERROR: Commit message length too long (limit is " + 70 + "): " + commit.subject)
    return 1

  if commit.style["type"] == "":
    sys.stderr.write("ERROR: Commit message must have a type 'type: subject': " + commit.subject)
    return 2

  if commit.style["type"] in ["Features", "Bug Fixes"]:
    if commit.style["scope"] is not None and JIRA_ISSUE_SYMBOL not in commit.style["scope"]:
      sys.stderr.write("ERROR: Scope for Feature and Bug fix needs to be Jira issue in format '"+JIRA_ISSUE_SYMBOL+"XXX': " + commit.style["scope"])
      return 2

    if "issues" not in commit.text_refs or len(commit.text_refs["issues"]) == 0:
      sys.stderr.write("ERROR: Feature and bug fix must have a Jira issue linked")
      sys.stderr.write("ERROR: You can link the issue either in subject as 'feat("+JIRA_ISSUE_SYMBOL+"XXX): subject': " + commit.subject)
      sys.stderr.write("ERROR: or anywhere in the body preferably as 'Fixes: "+JIRA_ISSUE_SYMBOL+"XXX' or Refs: "+JIRA_ISSUE_SYMBOL+"XXX")
      return 3

  return 0

if __name__ == '__main__':
  if sys.argv[1] == 'commit-check':
    sys.exit(check_commit())
  else:
    sys.exit(main(sys.argv[1:]))
