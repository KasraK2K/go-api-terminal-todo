Build Project
```bash
go build -o todo_cli cmd/cli.go
```

Move built application file and database file to your destination path
```bash
mv todo_cli <destination_folder>
mkdir <destination_folder>/db
mkdir <destination_folder>/db/database.db
```

Create alias in you favorite terminal (I have used zsh so place it into .zshrc)
```bash
todo() {
    current_dir="$PWD"  # Save the current directory
    cd <destination_folder> || return
    ./todo_cli "$@"      # Run the CLI with all arguments
    cd "$current_dir"    # Return to the original directory
}
```

Fix permissions
```bash
chmod -R 755 <destination_folder>
chmod 666 db/database.db
```

Test commands
```bash
todo create -t "My First Task" -d "This is my first task"
todo list
todo update -i 1 -c true
```