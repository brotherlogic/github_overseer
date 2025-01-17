# Repository Configuration

author:brotherlogic
updated:2025-01-06
type:feature

## Abstract

This module configures existing github repositories to support basic
functionality that we expect from our github repos. It supports a baseline level of support
and then specializations depending on the nature of the repository.

## Process

Runs in the background, every 5 minutes.

1. Gets a list of all the repos
1. Gets the .repository file, reads the text proto within
    1. If the .repository file doesn't exist, then raise an issue to add one
    1. Get the repository type:
        1. CODE
        1. PROCESS
    1. Configure the repo appropriately
        1. Adds in the PERSONAL_TOKEN secret

## Tasks

1. Overseer runs a prints the count of all the repos
1. githubridge: Support reading file from repo
1. Repository proto defined
1. Overseer reads .repository file and parses proto
1. Raise issue on missing repo file
1. githubridge: support setting secrets
1. If repo type is CODE -> add personal token secret
