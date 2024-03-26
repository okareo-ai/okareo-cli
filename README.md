
# Okareo CLI
# Continuous Model Delivery

A tool for interacting with the Okareo API

## USAGE
1. Download okareo from dist
2. Add okareo to your PATH (export PATH=$PATH:/correct/location/to/okareo)
3. In your repo, add a .okareo folder
4. Add .okareo/config.yml
5. Add your test scripts in .okareo/flows

- .okareo
    - config.yml
    - flows
        - classification.py
        - retrieval.py
        - generation.py
        - edit_distance.py

## Flags on Okareo
- '-latest=false' Okareo will install okareo on every run by default.  If you want to control version, set this to false.
- '-file=ABC' This allows you to use okareo to run a specific file
- '-debug' If you are unsure what is happening in okareo, you can use this to see more

## Okareo config.yml
```
name: <NAME OF YOUR REPO>
api-key: ${OKAREO_API_KEY}
project-id: ${OKAREO_PROJECT_ID}
run:
  scripts:
    file-pattern: '.*\.py'
```