# Tolite

[![Build Status](https://travis-ci.com/timedia/tolite.svg?branch=master)](https://travis-ci.com/timedia/tolite)

gitolite.conf を管理するツール

## TODO

サブコマンドの完成状況

- [x] update
- [ ] users
  - [ ] add
  - [ ] update
  - [ ] list
  - [ ] delete
- [ ] groups
  - [ ] add
  - [ ] update
  - [ ] list
  - [ ] delete
- [ ] repos
  - [ ] add
  - [ ] update
  - [ ] list
  - [ ] delete
- etc
  - [ ] config.yml

## 仕様

`tolite update` で自動的に gitolite.yml の変更がコミットされ、反映される。gitolite.yml がなく、gitolite.conf だけがある場合には gitolite.yml を新しく生成する。

### Users

`tolite users add <user name> [-e email] [-k key...] [-k key=public key...] [--key-file key=public key file...]`

`tolite users update [--update-name | -n new user name] [-e email] [--delete-email] [-k key...] [--delete-key key...] <user name>`

`--delete-key` はそのうちキーのファイルも連動して消えるようにする

`tolite users list`

`tolite users delete <user name>`

### Groups

`tolite groups add <group name> [users...]`

`tolite groups update <group name> [--update-name | -n new group name] [-u --add-user user name...] [-d --delete-user user name...]`

`tolite groups list`

`tolite groups delete <group name>`

### Repos

`tolite repos add <repo name> [--ro group...] [--rw group...] [--rw+ group...]`

`tolite repos update <repo name> [-d | --delete] <r | rw | rw+> group...`

`tolite repos list`

`tolite repos delete <repo name>`

Admin Repo も同様

### Config file

```gitolite.yml
users:
  username:
    email: user@timedia.co.jp
    keys: key_name
  another_user:
    keys:
      another_key_name
      superawesomekey: raw+public+key+string
groups:
  group_name:
    - username
    - another_user
  another_group:
    - username
    - hoge
    - fuga
repos:
  repo_dir:
    RW+:
      - group_name
      - another_group
    R:
      - group_name_readonly
admin_repos:
  admin_repo/dir:
    RW+: all
```

- keys は名前だけの場合は鍵のフォルダから引っ張ってくる。ハッシュの場合は鍵ファイルを `pubkey` 以下に生成する。
- user の email は任意とする。

```.config/tolite/config.yml
gitoliteYmlPath: ''
gitoliteConfPath: ''
pubkeyDirPath: ''
```

## 名前の由来

フランス語で TNT です
