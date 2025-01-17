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
