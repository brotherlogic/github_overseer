# Task tracker

## Background

This is a replicate of the task tracker in the old cluster. With a
new repo twist.

## Process

The tracker searches for markdown files in active repositories and then tracks them.
Specifically it tracks for certains classes of files:

1. files called _todo where we have structured lists of things
   that need to be done. These are files with potentially a preamble but only containing
   an ordered list, two levels deep. The first level is the class of todo, then second
   level are the specific items. These are tracked as tasklists under each pile
1. Non _todo files that contain a section called "Tasks". These are then created as tasklists.

Tasklists are created under the repo that created the files, unless tasks are preceeded by a name then
a colon, in which case they are placed under that repo.

## Tasks

1. githubridge: Githubridge supports listing files in the repo
1. Overseer runs through each project and builds a mapping from repo -> latest hash
1. Overseer tracks file listing to given hash
1. Overseer can scan a repo and find md files
1. githubridge: Githubridge supports listing contents of file in repo
1. Overseer can analyze a file and pull tasks out of a given file
1. Build a representation of the task state for the given hash
1. Find the first open task and file a bug and adjust the document to attach the bug ref
1. Overseer checks task bugs on startup - if issue is closed, add a strikeout to the document
