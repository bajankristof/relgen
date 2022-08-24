# RelGen
RelGen is a command-line tool that helps you automate version bumps and changelog creation for any of your projects by using Conventional Commits and your Git history.

```json
{
  "preRelease": "",
  "versionPrefix": false,
  "buildMetadata": "",
  "changeSpec": [{
    "type": "^feat$",
    "bump": "MINOR",
    "category": "Features"
  }, {
    "type": "^fix$",
    "bump": "PATCH",
    "category": "Fixes"
  },{
    "type": "^build|chore|ci|docs|style|refactor|perf|test$",
    "bump": "NONE",
    "category": "Other"
  }]
}
```
