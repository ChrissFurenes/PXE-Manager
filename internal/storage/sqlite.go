package storage

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Store struct {
	DB *sql.DB
}

func New(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	s := &Store{DB: db}
	if err := s.init(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) init() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			boot_mode TEXT NOT NULL,
			boot_type TEXT NOT NULL DEFAULT 'kernel_initrd',
			kernel TEXT NOT NULL DEFAULT '',
			initrd TEXT NOT NULL DEFAULT '',
			image_path TEXT NOT NULL DEFAULT '',
			cmdline TEXT NOT NULL DEFAULT '',
			enabled INTEGER NOT NULL DEFAULT 1
		);`,
		`CREATE TABLE IF NOT EXISTS clients (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			mac TEXT NOT NULL UNIQUE,
			hostname TEXT NOT NULL,
			profile_id INTEGER,
			config_id INTEGER,
			show_menu INTEGER NOT NULL DEFAULT 0,
			description TEXT NOT NULL DEFAULT '',
			FOREIGN KEY(profile_id) REFERENCES profiles(id),
    		FOREIGN KEY(config_id) REFERENCES talos(id)
		);`,
		`CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS assets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			file_name TEXT NOT NULL,
			file_path TEXT NOT NULL,
			file_type TEXT NOT NULL,
			size_bytes INTEGER NOT NULL DEFAULT 0,
			description TEXT NOT NULL DEFAULT ''
		);`,
		`CREATE TABLE IF NOT EXISTS talos (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		name TEXT NOT NULL,
    		file_name TEXT NOT NULL,
    		file_path TEXT NOT NULL,
    		size_bytes INTEGER NOT NULL DEFAULT 0,
    		description TEXT NOT NULL DEFAULT ''
		);`,
	}

	for _, stmt := range stmts {
		if _, err := s.DB.Exec(stmt); err != nil {
			return err
		}
	}

	if _, err := s.DB.Exec(`INSERT OR IGNORE INTO settings (key, value) VALUES ('default_profile_id', '0')`); err != nil {
		return err
	}

	return nil
}

func (s *Store) ListProfiles() ([]Profile, error) {
	rows, err := s.DB.Query(`SELECT id, name, boot_mode, boot_type, kernel, initrd, image_path, cmdline, enabled FROM profiles ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Profile
	for rows.Next() {
		var p Profile
		var enabled int
		if err := rows.Scan(&p.ID, &p.Name, &p.BootMode, &p.BootType, &p.Kernel, &p.Initrd, &p.ImagePath, &p.Cmdline, &enabled); err != nil {
			return nil, err
		}
		p.Enabled = enabled == 1
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) GetProfile(id int64) (*Profile, error) {
	row := s.DB.QueryRow(`SELECT id, name, boot_mode, boot_type, kernel, initrd, image_path, cmdline, enabled FROM profiles WHERE id = ?`, id)

	var p Profile
	var enabled int
	if err := row.Scan(&p.ID, &p.Name, &p.BootMode, &p.BootType, &p.Kernel, &p.Initrd, &p.ImagePath, &p.Cmdline, &enabled); err != nil {
		return nil, err
	}
	p.Enabled = enabled == 1
	return &p, nil
}

func (s *Store) CreateProfile(p *Profile) error {
	enabled := 0
	if p.Enabled {
		enabled = 1
	}

	res, err := s.DB.Exec(
		`INSERT INTO profiles (name, boot_mode, boot_type, kernel, initrd, image_path, cmdline, enabled)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Name, p.BootMode, p.BootType, p.Kernel, p.Initrd, p.ImagePath, p.Cmdline, enabled,
	)
	if err != nil {
		return err
	}
	p.ID, _ = res.LastInsertId()
	return nil
}

func (s *Store) UpdateProfile(p *Profile) error {
	enabled := 0
	if p.Enabled {
		enabled = 1
	}

	_, err := s.DB.Exec(
		`UPDATE profiles
		 SET name = ?, boot_mode = ?, boot_type = ?, kernel = ?, initrd = ?, image_path = ?, cmdline = ?, enabled = ?
		 WHERE id = ?`,
		p.Name, p.BootMode, p.BootType, p.Kernel, p.Initrd, p.ImagePath, p.Cmdline, enabled, p.ID,
	)
	return err
}

func (s *Store) DeleteProfile(id int64) error {
	_, err := s.DB.Exec(`DELETE FROM profiles WHERE id = ?`, id)
	return err
}

func (s *Store) ListClients() ([]Client, error) {
	rows, err := s.DB.Query(`SELECT id, mac, hostname, profile_id, config_id, show_menu, description FROM clients ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Client
	for rows.Next() {
		var c Client
		var show int
		if err := rows.Scan(&c.ID, &c.MAC, &c.Hostname, &c.ProfileID, &c.ConfigID, &show, &c.Description); err != nil {
			return nil, err
		}
		c.ShowMenu = show == 1
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *Store) GetClient(id int64) (*Client, error) {
	row := s.DB.QueryRow(`SELECT id, mac, hostname, profile_id, config_id, show_menu, description FROM clients WHERE id = ?`, id)

	var c Client
	var show int
	if err := row.Scan(&c.ID, &c.MAC, &c.Hostname, &c.ProfileID, &c.ConfigID, &show, &c.Description); err != nil {
		return nil, err
	}
	c.ShowMenu = show == 1
	return &c, nil
}

func (s *Store) GetClientByMAC(mac string) (*Client, error) {
	row := s.DB.QueryRow(`SELECT id, mac, hostname, profile_id, config_id, show_menu, description FROM clients WHERE lower(mac) = lower(?)`, mac)

	var c Client
	var show int
	if err := row.Scan(&c.ID, &c.MAC, &c.Hostname, &c.ProfileID, &c.ConfigID, &show, &c.Description); err != nil {
		return nil, err
	}
	c.ShowMenu = show == 1
	return &c, nil
}

func (s *Store) CreateClient(c *Client) error {
	show := 0
	if c.ShowMenu {
		show = 1
	}

	res, err := s.DB.Exec(
		`INSERT INTO clients (mac, hostname, profile_id, config_id, show_menu, description) VALUES (?, ?, ?, ?, ?, ?)`,
		c.MAC, c.Hostname, c.ProfileID, c.ConfigID, show, c.Description,
	)
	if err != nil {
		return err
	}
	c.ID, _ = res.LastInsertId()
	return nil
}

func (s *Store) UpdateClient(c *Client) error {
	show := 0
	if c.ShowMenu {
		show = 1
	}

	_, err := s.DB.Exec(
		`UPDATE clients
		 SET mac = ?, hostname = ?, profile_id = ?, config_id = ?, show_menu = ?, description = ?
		 WHERE id = ?`,
		c.MAC, c.Hostname, c.ProfileID, c.ConfigID, show, c.Description, c.ID,
	)
	return err
}

func (s *Store) DeleteClient(id int64) error {
	_, err := s.DB.Exec(`DELETE FROM clients WHERE id = ?`, id)
	return err
}

func (s *Store) GetDefaultProfileID() (int64, error) {
	row := s.DB.QueryRow(`SELECT value FROM settings WHERE key = 'default_profile_id'`)
	var raw string
	if err := row.Scan(&raw); err != nil {
		return 0, err
	}

	var id int64
	_, err := fmt.Sscanf(raw, "%d", &id)
	return id, err
}

func (s *Store) SetDefaultProfileID(id int64) error {
	_, err := s.DB.Exec(`INSERT INTO settings (key, value) VALUES ('default_profile_id', ?)
		ON CONFLICT(key) DO UPDATE SET value=excluded.value`, fmt.Sprintf("%d", id))
	return err
}

func (s *Store) ListAssets() ([]Asset, error) {
	rows, err := s.DB.Query(`SELECT id, name, file_name, file_path, file_type, size_bytes, description FROM assets ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Asset
	for rows.Next() {
		var a Asset
		if err := rows.Scan(&a.ID, &a.Name, &a.FileName, &a.FilePath, &a.FileType, &a.SizeBytes, &a.Description); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *Store) GetAsset(id int64) (*Asset, error) {
	row := s.DB.QueryRow(`SELECT id, name, file_name, file_path, file_type, size_bytes, description FROM assets WHERE id = ?`, id)

	var a Asset
	if err := row.Scan(&a.ID, &a.Name, &a.FileName, &a.FilePath, &a.FileType, &a.SizeBytes, &a.Description); err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Store) CreateAsset(a *Asset) error {
	res, err := s.DB.Exec(
		`INSERT INTO assets (name, file_name, file_path, file_type, size_bytes, description)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		a.Name, a.FileName, a.FilePath, a.FileType, a.SizeBytes, a.Description,
	)
	if err != nil {
		return err
	}
	a.ID, _ = res.LastInsertId()
	return nil
}

func (s *Store) UpdateAsset(a *Asset) error {
	_, err := s.DB.Exec(
		`UPDATE assets
		 SET name = ?, file_type = ?, description = ?
		 WHERE id = ?`,
		a.Name, a.FileType, a.Description, a.ID,
	)
	return err
}

func (s *Store) DeleteAsset(id int64) error {
	_, err := s.DB.Exec(`DELETE FROM assets WHERE id = ?`, id)
	return err
}

func (s *Store) ListTalos() ([]Talos, error) {
	rows, err := s.DB.Query(`SELECT id, name, file_name, file_path, size_bytes, description FROM talos ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Talos
	for rows.Next() {
		var c Talos
		if err := rows.Scan(&c.ID, &c.Name, &c.FileName, &c.FilePath, &c.SizeBytes, &c.Description); err != nil {
			return nil, err
		}

		out = append(out, c)
	}
	return out, rows.Err()
}
func (s *Store) CreateTalos(a *Talos) error {

	res, err := s.DB.Exec(
		`INSERT INTO talos (name, file_name, file_path, size_bytes, description)
		 VALUES (?, ?, ?, ?, ?)`,
		a.Name, a.FileName, a.FilePath, a.SizeBytes, a.Description,
	)
	if err != nil {
		return err
	}
	a.ID, _ = res.LastInsertId()
	return nil
}
func (s *Store) GetTalos(id int64) (*Talos, error) {
	row := s.DB.QueryRow(`SELECT id, name, file_name, file_path, size_bytes, description FROM talos WHERE id = ?`, id)

	var a Talos
	if err := row.Scan(&a.ID, &a.Name, &a.FileName, &a.FilePath, &a.SizeBytes, &a.Description); err != nil {
		return nil, err
	}
	return &a, nil
}
func (s *Store) UpdateTalos(a *Talos) error {
	_, err := s.DB.Exec(
		`UPDATE talos
		 SET name = ?, description = ?
		 WHERE id = ?`,
		a.Name, a.Description, a.ID,
	)
	return err
}

func (s *Store) DeleteTalos(id int64) error {
	_, err := s.DB.Exec(`DELETE FROM talos WHERE id = ?`, id)
	return err
}
