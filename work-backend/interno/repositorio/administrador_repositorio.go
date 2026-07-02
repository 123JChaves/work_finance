package repositorio

import (
	"errors"
	"sync"
	"work-backend/interno/entidade"
)

type MemoriaAdministradorRepositorio struct {
	sync.RWMutex
	administradores map[int]entidade.Administrador
	proximoID       int
}

func NovoMemoriaAdministradorRepositorio() *MemoriaAdministradorRepositorio {
	return &MemoriaAdministradorRepositorio{
		administradores: make(map[int]entidade.Administrador),
		proximoID:       1,
	}
}

func (repositorio *MemoriaAdministradorRepositorio) Create(administrador *entidade.Administrador) error {
	repositorio.Lock()
	defer repositorio.Unlock()

	administrador.ID = repositorio.proximoID
	repositorio.proximoID++
	repositorio.administradores[administrador.ID] = *administrador
	return nil
}

func (repositorio *MemoriaAdministradorRepositorio) FindByID(id int) (*entidade.Administrador, error) {
	repositorio.RLock()
	defer repositorio.RUnlock()

	administrador, existe := repositorio.administradores[id]
	if !existe {
		return nil, errors.New("Administrador não encontrado")
	}
	return &administrador, nil
}

func (repositorio *MemoriaAdministradorRepositorio) FindAll() ([]*entidade.Administrador, error) {
	repositorio.RLock()
	defer repositorio.RUnlock()

	lista := make([]*entidade.Administrador, 0, len(repositorio.administradores))
	for _, admin := range repositorio.administradores {
		// Criamos uma cópia para evitar ponteiros compartilhados indevidamente
		copia := admin
		lista = append(lista, &copia)
	}

	return lista, nil
}


func (repositorio *MemoriaAdministradorRepositorio) FindByEmail(email string) (*entidade.Administrador, error) {
	repositorio.RLock()
	defer repositorio.RUnlock()

	for _, administrador := range repositorio.administradores {
		if administrador.Email == email {
			return &administrador, nil
		}
	}
	return nil, nil 
}

func (repositorio *MemoriaAdministradorRepositorio) Update(administrador *entidade.Administrador) error {
	repositorio.Lock()
	defer repositorio.Unlock() // Corrigido de r.Unlock() para repositorio.Unlock()

	_, existe := repositorio.administradores[administrador.ID]
	if !existe {
		return errors.New("Administrador não encontrado para atualização")
	}
	repositorio.administradores[administrador.ID] = *administrador
	return nil
}

func (repositorio *MemoriaAdministradorRepositorio) Delete(id int) error {
	repositorio.Lock()
	defer repositorio.Unlock()

	_, existe := repositorio.administradores[id] // Corrigido de r.administradores para repositorio.administradores
	if !existe {
		return errors.New("Administrador não encontrado para exclusão")
	}
	delete(repositorio.administradores, id)
	return nil
}