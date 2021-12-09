Just Another Build System
=========================


TODO
----

- [ ] Defaults conventions
- [ ] Timestamps bank
- [ ] Standard library
- [ ] Other formats output
- [ ] Output transformers
- [ ] Build server


Feature list
------------

- Instructions consist of macros and rules
- Macros:
	- Macros can have explicit values, or can be shell command results
	- Shell commands can either be shell cmdlets, or shebang scriptlets
- Rules:
	- Can depend on other rules
	- Can be conditionally executed by observing shellcommand status: non-zero exit code triggers the rule execution
	- Can have series of tasks as shellcommands
- Logging with levels
- Tasks can be either printed in shell script compatible way, or executed
- Has dedicated subcommands for printing and running and their own cli flags
- Is able to monitor the list of files in a directory descruibed by `--dir` and `-filter` flags via the `watch` command, and re-run on change in mode designated by the `--action` param: `go run . watch --action=run cover`
- Action/subcommand outputs are sent to a channel, separate debug/log info and standard out
- Option to stop and plow through on execution error (`--force`)
- Can optionally include rule conditions in print output
- Can toggle verbosity from CLI params
- Can load multi-file definitions by making use of file globs
