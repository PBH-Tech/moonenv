---
name: Bug
description: Create a report to help us improve Moonenv by reproducing and fixing bugs
title: "[Bug]: "
labels:
  - "bug_report"

body:
  - type: markdown
    attributes:
      value: |
        **Thanks for taking the time to fill out this bug report!**

  - type: textarea
    attributes:
      label: Current Behavior
      description: |
        A clear and concise description of the issue you're experiencing.
    validations:
      required: false

  - type: textarea
    attributes:
      label: Expected Behavior
      description: |
        A clear and concise description of what you expected to happen.
    validations:
      required: false

  - type: textarea
    attributes:
      label: Steps To Reproduce
      description: Steps to reproduce the behavior.
    validations:
      required: true

  - type: textarea
    attributes:
      label: Screenshots
      description: |
        If applicable, add screenshots to help explain your problem.

        Tip: You can attach images or log files by clicking this area to highlight it and then dragging files in.
    validations:
      required: false

  - type: textarea
    attributes:
      label: Environment
      description: |
        Information about your environment where you see the bug.

      value: |
        OS:

        Version:

    validations:
      required: true

  - type: textarea
    attributes:
      label: Additional Context
      description: |
        Add any other context about the problem here.
    validations:
      required: false