Just Another Build System
=========================


TODO
----

- [ ] Other formats output
- [ ] Output transformers
- [ ] Build server
- [ ] Defaults conventions
- [ ] Timestamps bank
- [ ] Standard library
- [ ] Multi-file definitions


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
- Is able to monitor the list of files supplied via stdin using `watch` command, and re-run on change in mode designated by the `--action` param: `find . -type f -name '*.go' | go run . watch --action=run cover`
- Action/subcommand outputs are sent to a channel, separate debug/log info and standard out
- Option to stop and plow through on execution error (`--force`)
