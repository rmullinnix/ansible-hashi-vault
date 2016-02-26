# ansible-hashi-vault

Ansible script to deploy hashicorp vault to existing consul cluster.  This is just for example purposes.  Feel free to copy and modify.

Usage: ansible-playbook -i hosts.dev vault.yml --ask-become-pass

Configuration
*  populate hosts.dev with servers (consulservers - hosts running consul in server mode, consulagents - hosts running consul in agent mode, vaultmaster - single server to invoke initialization)
*  Alter group_vars/development/all with your environment specific information
*  Alter group_vars/development/consul with consul specific information
*  Alter group_vars/development/vault with vault specific information
*  Build vault_init - cd roles/vault/files/vault_init; go build - copy to {{ build_dir }}/dev/vault/
*  Copy vault executable to {{ build_dir }}/dev/vault/

Issues
* Some of the tasks are linux specifc commands (SUSE for me) and may have to be altered
* During the roles/vault/tasks/init.yml execution, in my enviornment vault receives a SIGPIPE and exits.  I had to alter command/server.go line 105 to use os.Stdout as opposed to os.Stderr for it to execute successfully.
* Vault keys are pulled back and stored in roles/vault/vars - be sure to secure it after the run
