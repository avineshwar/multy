CHANGELOG
=========

## [0.1.3-alpha] - 16/01/2021

### Changed

Changed virtual_machine resource references. Now passing full object instead of id

---

## [0.1.2-alpha] - 15/01/2021

### Added

`object_storage_object.acl` parameter

---

## [0.1.1-alpha] - 14/01/2021

### Changed

CLI Logging

### Fixed

`virtual_network.user_data` improvements

---

## [0.1.0-alpha] - 10/01/2021

### Added

`virtual_machine` - added `ssh_key_file_path` parameter

### Changed

[Azure] `virtual_machine` - password now uses terraform `random_password` instead of hardcoded password

---

## [0.0.1-alpha] - 08/01/2021

Initial release
