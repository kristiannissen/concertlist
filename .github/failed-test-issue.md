---
title: "Tests failing: {{ env.GITHUB_REF_NAME }}"
labels: bug
---

Workflow **{{ env.GITHUB_WORKFLOW }}** failed run **#{{ env.GITHUB_RUN_NUMBER }}**.

* **Commit:** {{ env.GITHUB_SHA }}
* **Details:** https://github.com/{{ env.GITHUB_REPOSITORY }}/actions/runs/{{ env.GITHUB_RUN_ID }}
