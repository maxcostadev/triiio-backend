package imoveis

import (
	"context"
	"errors"
	"fmt"
)

// Service defines the interface for property business logic
type Service interface {
	// Imovel Operations
	CreateImovel(ctx context.Context, req *CreateImovelRequest) (*ImovelResponse, error)
	GetImovel(ctx context.Context, id uint) (*ImovelResponse, error)
	GetImovelByCodigo(ctx context.Context, codigo string) (*ImovelResponse, error)
	GetImovelByIdIntegracao(ctx context.Context, idIntegracao string) (*ImovelResponse, error)
	UpdateImovel(ctx context.Context, id uint, req *UpdateImovelRequest) (*ImovelResponse, error)
	DeleteImovel(ctx context.Context, id uint) error
	HardDeleteImovel(ctx context.Context, id uint) error

	// List & Filter
	ListImoveis(ctx context.Context, query *ImovelListQuery) (*ImovelListResponse, error)
	ListImovelsByEmpreendimento(ctx context.Context, empreendimentoID uint, page, limit int) ([]ImovelResponse, int64, error)
	ListImovelsByOrganizacao(ctx context.Context, organizacaoID uint, page, limit int) ([]ImovelResponse, int64, error)

	// Bulk Operations
	CreateImovelBatch(ctx context.Context, reqs []CreateImovelRequest) error
	UpdateImovelBatch(ctx context.Context, imoveis []Imovel) error

	// Statistics
	CountImoveis(ctx context.Context) (int64, error)
	CountImovelsByStatus(ctx context.Context, status string) (int64, error)
	CountImovelsByEmpreendimento(ctx context.Context, empreendimentoID uint) (int64, error)

	// Existence checks
	ImovelExistsByCodigo(ctx context.Context, codigo string) (bool, error)
	ImovelExistsByIdIntegracao(ctx context.Context, idIntegracao string) (bool, error)

	// Relationship Operations - Anexos
	AddAnexo(ctx context.Context, imovelID uint, anexo *Anexo) error
	RemoveAnexo(ctx context.Context, imovelID, anexoID uint) error
	GetAnexos(ctx context.Context, imovelID uint) ([]AnexoResponse, error)

	// Relationship Operations - Single associations
	AttachEndereco(ctx context.Context, imovelID, enderecoID uint) error
	AttachEmpreendimento(ctx context.Context, imovelID, empreendimentoID uint) error
	AttachPlanta(ctx context.Context, imovelID, plantaID uint) error
	AttachPacote(ctx context.Context, imovelID, pacoteID uint) error
	AttachOrganizacao(ctx context.Context, imovelID, organizacaoID uint) error
	AttachPrecoVenda(ctx context.Context, imovelID, precoVendaID uint) error
	AttachPrecoAluguel(ctx context.Context, imovelID, precoAluguelID uint) error

	// Endereco Operations (for import/external integration)
	CreateEndereco(ctx context.Context, endereco *Endereco) error

	// Relationship Operations - Caracteristicas
	AddCaracteristicas(ctx context.Context, imovelID uint, caracteristicaIDs []uint) error
	RemoveCaracteristicas(ctx context.Context, imovelID uint, caracteristicaIDs []uint) error
	GetCaracteristicas(ctx context.Context, imovelID uint) ([]CaracteristicaResponse, error)
	ReplaceCaracteristicas(ctx context.Context, imovelID uint, caracteristicaIDs []uint) error
}

type service struct {
	repo Repository
}

// NewService creates a new property service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateImovel creates a new property
func (s *service) CreateImovel(ctx context.Context, req *CreateImovelRequest) (*ImovelResponse, error) {
	// Validate business rules
	if req.Objetivo == "ALUGAR" && req.PrecoAluguelID == 0 {
		return nil, fmt.Errorf("rental properties must have a rental price")
	}
	if req.Objetivo == "VENDER" && req.PrecoVendaID == 0 {
		return nil, fmt.Errorf("properties for sale must have a selling price")
	}

	// Check if codigo already exists
	exists, err := s.repo.ExistsByCodigo(ctx, req.Codigo)
	if err != nil {
		return nil, fmt.Errorf("failed to check codigo uniqueness: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("property with codigo '%s' already exists", req.Codigo)
	}

	// Check if idIntegracao is unique (if provided)
	if req.IdIntegracao != "" {
		exists, err := s.repo.ExistsByIdIntegracao(ctx, req.IdIntegracao)
		if err != nil {
			return nil, fmt.Errorf("failed to check idIntegracao uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("property with idIntegracao '%s' already exists", req.IdIntegracao)
		}
	}

	// Create model from request
	imovel := &Imovel{
		Id_Integracao:       req.IdIntegracao,
		Titulo:              req.Titulo,
		Codigo:              req.Codigo,
		Tipo:                req.Tipo,
		Objetivo:            req.Objetivo,
		Finalidade:          req.Finalidade,
		Descricao:           req.Descricao,
		Metragem:            req.Metragem,
		NumQuartos:          req.NumQuartos,
		NumSuites:           req.NumSuites,
		NumBanheiros:        req.NumBanheiros,
		NumVagas:            req.NumVagas,
		NumAndar:            req.NumAndar,
		Unidade:             req.Unidade,
		Condominio:          req.Condominio,
		IPTU:                req.IPTU,
		InscricaoIPTU:       req.InscricaoIPTU,
		EnderecoID:          req.EnderecoID,
		PlantaID:            req.PlantaID,
		CorretorPrincipalID: req.CorretorPrincipalID,
		PacoteID:            req.PacoteID,
		Status:              "EM_EDICAO", // Default status
		Published:           false,
		Closed:              false,
	}

	// Only set optional foreign keys if they're provided (non-zero)
	if req.EmpreendimentoID != 0 {
		imovel.EmpreendimentoID = req.EmpreendimentoID
	}
	if req.PrecoVendaID != 0 {
		imovel.PrecoVendaID = req.PrecoVendaID
	}
	if req.PrecoAluguelID != 0 {
		imovel.PrecoAluguelID = req.PrecoAluguelID
	}

	// Build list of fields to omit based on zero values
	omitFields := []string{}
	if req.EmpreendimentoID == 0 {
		omitFields = append(omitFields, "EmpreendimentoID")
	}
	if req.PrecoVendaID == 0 {
		omitFields = append(omitFields, "PrecoVendaID")
	}
	if req.PrecoAluguelID == 0 {
		omitFields = append(omitFields, "PrecoAluguelID")
	}
	if req.PlantaID == 0 {
		omitFields = append(omitFields, "PlantaID")
	}
	if req.CorretorPrincipalID == 0 {
		omitFields = append(omitFields, "CorretorPrincipalID")
	}
	if req.PacoteID == 0 {
		omitFields = append(omitFields, "PacoteID")
	}

	// Save to repository with omitted fields
	if err := s.repo.(*repository).db.Omit(omitFields...).Create(imovel).Error; err != nil {
		return nil, fmt.Errorf("failed to create property: %w", err)
	}

	// Retrieve and return
	return s.GetImovel(ctx, imovel.ID)
}

// GetImovel retrieves a property by ID
func (s *service) GetImovel(ctx context.Context, id uint) (*ImovelResponse, error) {
	if id == 0 {
		return nil, errors.New("invalid property ID")
	}

	imovel, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve property: %w", err)
	}

	if imovel == nil {
		return nil, fmt.Errorf("property not found")
	}

	return s.mapToResponse(imovel), nil
}

// GetImovelByCodigo retrieves a property by codigo
func (s *service) GetImovelByCodigo(ctx context.Context, codigo string) (*ImovelResponse, error) {
	if codigo == "" {
		return nil, errors.New("codigo cannot be empty")
	}

	imovel, err := s.repo.FindByCodigo(ctx, codigo)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve property: %w", err)
	}

	if imovel == nil {
		return nil, fmt.Errorf("property with codigo '%s' not found", codigo)
	}

	return s.mapToResponse(imovel), nil
}

// GetImovelByIdIntegracao retrieves a property by integration ID
func (s *service) GetImovelByIdIntegracao(ctx context.Context, idIntegracao string) (*ImovelResponse, error) {
	if idIntegracao == "" {
		return nil, errors.New("idIntegracao cannot be empty")
	}

	imovel, err := s.repo.FindByIdIntegracao(ctx, idIntegracao)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve property: %w", err)
	}

	if imovel == nil {
		return nil, fmt.Errorf("property with idIntegracao '%s' not found", idIntegracao)
	}

	return s.mapToResponse(imovel), nil
}

// UpdateImovel updates an existing property
func (s *service) UpdateImovel(ctx context.Context, id uint, req *UpdateImovelRequest) (*ImovelResponse, error) {
	if id == 0 {
		return nil, errors.New("invalid property ID")
	}

	// Get existing property
	imovel, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve property: %w", err)
	}

	if imovel == nil {
		return nil, fmt.Errorf("property not found")
	}

	// Check for codigo uniqueness if changing it
	if req.Codigo != "" && req.Codigo != imovel.Codigo {
		exists, err := s.repo.ExistsByCodigo(ctx, req.Codigo)
		if err != nil {
			return nil, fmt.Errorf("failed to check codigo uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("property with codigo '%s' already exists", req.Codigo)
		}
		imovel.Codigo = req.Codigo
	}

	// Update fields
	if req.Titulo != "" {
		imovel.Titulo = req.Titulo
	}
	if req.Tipo != "" {
		imovel.Tipo = req.Tipo
	}
	if req.Objetivo != "" {
		imovel.Objetivo = req.Objetivo
	}
	if req.Finalidade != "" {
		imovel.Finalidade = req.Finalidade
	}
	if req.Descricao != "" {
		imovel.Descricao = req.Descricao
	}
	if req.Metragem != nil && *req.Metragem > 0 {
		imovel.Metragem = *req.Metragem
	}
	if req.NumQuartos != nil && *req.NumQuartos >= 0 {
		imovel.NumQuartos = *req.NumQuartos
	}
	if req.NumSuites != nil && *req.NumSuites >= 0 {
		imovel.NumSuites = *req.NumSuites
	}
	if req.NumBanheiros != nil && *req.NumBanheiros >= 0 {
		imovel.NumBanheiros = *req.NumBanheiros
	}
	if req.NumVagas != nil && *req.NumVagas >= 0 {
		imovel.NumVagas = *req.NumVagas
	}
	if req.NumAndar != nil {
		imovel.NumAndar = *req.NumAndar
	}
	if req.Unidade != "" {
		imovel.Unidade = req.Unidade
	}
	if req.Condominio != nil && *req.Condominio >= 0 {
		imovel.Condominio = *req.Condominio
	}
	if req.IPTU != nil && *req.IPTU >= 0 {
		imovel.IPTU = *req.IPTU
	}
	if req.InscricaoIPTU != "" {
		imovel.InscricaoIPTU = req.InscricaoIPTU
	}

	// Update relationships if provided
	if req.EnderecoID != nil {
		imovel.EnderecoID = *req.EnderecoID
	}
	if req.EmpreendimentoID != nil {
		imovel.EmpreendimentoID = *req.EmpreendimentoID
	}
	if req.PlantaID != nil {
		imovel.PlantaID = *req.PlantaID
	}
	if req.CorretorPrincipalID != nil {
		imovel.CorretorPrincipalID = *req.CorretorPrincipalID
	}
	if req.PacoteID != nil {
		imovel.PacoteID = *req.PacoteID
	}
	if req.PrecoVendaID != nil {
		imovel.PrecoVendaID = *req.PrecoVendaID
	}
	if req.PrecoAluguelID != nil {
		imovel.PrecoAluguelID = *req.PrecoAluguelID
	}

	// Update status fields
	if req.Status != "" {
		imovel.Status = req.Status
	}
	if req.Published != nil {
		imovel.Published = *req.Published
	}
	if req.Closed != nil {
		imovel.Closed = *req.Closed
	}

	// Update in repository
	if err := s.repo.Update(ctx, imovel); err != nil {
		return nil, fmt.Errorf("failed to update property: %w", err)
	}

	// Retrieve and return updated property
	return s.GetImovel(ctx, id)
}

// DeleteImovel soft deletes a property
func (s *service) DeleteImovel(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid property ID")
	}

	// Verify property exists
	imovel, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to verify property: %w", err)
	}

	if imovel == nil {
		return fmt.Errorf("property not found")
	}

	// Soft delete
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete property: %w", err)
	}

	return nil
}

// HardDeleteImovel permanently deletes a property
func (s *service) HardDeleteImovel(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid property ID")
	}

	// Verify property exists
	imovel, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to verify property: %w", err)
	}

	if imovel == nil {
		return fmt.Errorf("property not found")
	}

	// Hard delete
	if err := s.repo.HardDelete(ctx, id); err != nil {
		return fmt.Errorf("failed to permanently delete property: %w", err)
	}

	return nil
}

// ListImoveis retrieves properties with filtering and pagination
func (s *service) ListImoveis(ctx context.Context, query *ImovelListQuery) (*ImovelListResponse, error) {
	// Validate pagination parameters
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	// Retrieve from repository
	result, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list properties: %w", err)
	}

	return result, nil
}

// ListImovelsByEmpreendimento retrieves properties by enterprise
func (s *service) ListImovelsByEmpreendimento(ctx context.Context, empreendimentoID uint, page, limit int) ([]ImovelResponse, int64, error) {
	if empreendimentoID == 0 {
		return nil, 0, errors.New("invalid enterprise ID")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Retrieve from repository
	imoveis, total, err := s.repo.ListByEmpreendimento(ctx, empreendimentoID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list properties by enterprise: %w", err)
	}

	// Convert to responses
	responses := make([]ImovelResponse, len(imoveis))
	for i := range imoveis {
		responses[i] = *s.mapToResponse(&imoveis[i])
	}

	return responses, total, nil
}

// ListImovelsByOrganizacao retrieves properties by organization
func (s *service) ListImovelsByOrganizacao(ctx context.Context, organizacaoID uint, page, limit int) ([]ImovelResponse, int64, error) {
	if organizacaoID == 0 {
		return nil, 0, errors.New("invalid organization ID")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Retrieve from repository
	imoveis, total, err := s.repo.ListByCorretorPrincipal(ctx, organizacaoID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list properties by organization: %w", err)
	}

	// Convert to responses
	responses := make([]ImovelResponse, len(imoveis))
	for i := range imoveis {
		responses[i] = *s.mapToResponse(&imoveis[i])
	}

	return responses, total, nil
}

// CreateImovelBatch creates multiple properties
func (s *service) CreateImovelBatch(ctx context.Context, reqs []CreateImovelRequest) error {
	if len(reqs) == 0 {
		return errors.New("at least one property is required")
	}

	// Convert requests to models
	imoveis := make([]Imovel, len(reqs))
	for i, req := range reqs {
		// Basic validation per item
		if req.Codigo == "" {
			return fmt.Errorf("property at index %d: codigo is required", i)
		}

		imoveis[i] = Imovel{
			Id_Integracao:       req.IdIntegracao,
			Titulo:              req.Titulo,
			Codigo:              req.Codigo,
			Tipo:                req.Tipo,
			Objetivo:            req.Objetivo,
			Finalidade:          req.Finalidade,
			Descricao:           req.Descricao,
			Metragem:            req.Metragem,
			NumQuartos:          req.NumQuartos,
			NumSuites:           req.NumSuites,
			NumBanheiros:        req.NumBanheiros,
			NumVagas:            req.NumVagas,
			NumAndar:            req.NumAndar,
			Unidade:             req.Unidade,
			Condominio:          req.Condominio,
			IPTU:                req.IPTU,
			InscricaoIPTU:       req.InscricaoIPTU,
			EnderecoID:          req.EnderecoID,
			EmpreendimentoID:    req.EmpreendimentoID,
			PlantaID:            req.PlantaID,
			CorretorPrincipalID: req.CorretorPrincipalID,
			PacoteID:            req.PacoteID,
			PrecoVendaID:        req.PrecoVendaID,
			PrecoAluguelID:      req.PrecoAluguelID,
			Status:              "EM_EDICAO",
			Published:           false,
			Closed:              false,
		}
	}

	// Create batch in repository
	if err := s.repo.CreateBatch(ctx, imoveis); err != nil {
		return fmt.Errorf("failed to create properties in batch: %w", err)
	}

	return nil
}

// UpdateImovelBatch updates multiple properties
func (s *service) UpdateImovelBatch(ctx context.Context, imoveis []Imovel) error {
	if len(imoveis) == 0 {
		return errors.New("at least one property is required")
	}

	// Update batch in repository
	if err := s.repo.UpdateBatch(ctx, imoveis); err != nil {
		return fmt.Errorf("failed to update properties in batch: %w", err)
	}

	return nil
}

// CountImoveis returns total count of properties
func (s *service) CountImoveis(ctx context.Context) (int64, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count properties: %w", err)
	}
	return count, nil
}

// CountImovelsByStatus returns count of properties by status
func (s *service) CountImovelsByStatus(ctx context.Context, status string) (int64, error) {
	if status == "" {
		return 0, errors.New("status cannot be empty")
	}

	count, err := s.repo.CountByStatus(ctx, status)
	if err != nil {
		return 0, fmt.Errorf("failed to count properties by status: %w", err)
	}
	return count, nil
}

// CountImovelsByEmpreendimento returns count of properties by enterprise
func (s *service) CountImovelsByEmpreendimento(ctx context.Context, empreendimentoID uint) (int64, error) {
	if empreendimentoID == 0 {
		return 0, errors.New("invalid enterprise ID")
	}

	count, err := s.repo.CountByEmpreendimento(ctx, empreendimentoID)
	if err != nil {
		return 0, fmt.Errorf("failed to count properties by enterprise: %w", err)
	}
	return count, nil
}

// ImovelExistsByCodigo checks if a property exists by codigo
func (s *service) ImovelExistsByCodigo(ctx context.Context, codigo string) (bool, error) {
	if codigo == "" {
		return false, errors.New("codigo cannot be empty")
	}

	exists, err := s.repo.ExistsByCodigo(ctx, codigo)
	if err != nil {
		return false, fmt.Errorf("failed to check property existence: %w", err)
	}
	return exists, nil
}

// ImovelExistsByIdIntegracao checks if a property exists by integration ID
func (s *service) ImovelExistsByIdIntegracao(ctx context.Context, idIntegracao string) (bool, error) {
	if idIntegracao == "" {
		return false, errors.New("idIntegracao cannot be empty")
	}

	exists, err := s.repo.ExistsByIdIntegracao(ctx, idIntegracao)
	if err != nil {
		return false, fmt.Errorf("failed to check property existence: %w", err)
	}
	return exists, nil
}

// mapToResponse converts Imovel model to response DTO
func (s *service) mapToResponse(imovel *Imovel) *ImovelResponse {
	response := &ImovelResponse{
		ID:            imovel.ID,
		IdIntegracao:  imovel.Id_Integracao,
		Titulo:        imovel.Titulo,
		Codigo:        imovel.Codigo,
		SeqCodigo:     imovel.SeqCodigo,
		Tipo:          imovel.Tipo,
		Objetivo:      imovel.Objetivo,
		Finalidade:    imovel.Finalidade,
		Descricao:     imovel.Descricao,
		Metragem:      imovel.Metragem,
		NumQuartos:    imovel.NumQuartos,
		NumSuites:     imovel.NumSuites,
		NumBanheiros:  imovel.NumBanheiros,
		NumVagas:      imovel.NumVagas,
		NumAndar:      imovel.NumAndar,
		Unidade:       imovel.Unidade,
		Condominio:    imovel.Condominio,
		IPTU:          imovel.IPTU,
		InscricaoIPTU: imovel.InscricaoIPTU,
		Status:        imovel.Status,
		Published:     imovel.Published,
		Closed:        imovel.Closed,
		Visualizacoes: imovel.Visualizacoes,
		CreatedAt:     imovel.CreatedAt,
		UpdatedAt:     imovel.UpdatedAt,
	}

	// Map relationships
	if imovel.Endereco != nil {
		response.Endereco = &EnderecoResponse{
			ID:        imovel.Endereco.ID,
			Rua:       imovel.Endereco.Rua,
			Numero:    imovel.Endereco.Numero,
			Bairro:    imovel.Endereco.Bairro,
			Cidade:    imovel.Endereco.Cidade,
			Estado:    imovel.Endereco.Estado,
			CEP:       imovel.Endereco.CEP,
			Latitude:  imovel.Endereco.Latitude,
			Longitude: imovel.Endereco.Longitude,
		}
	}

	if imovel.Empreendimento != nil {
		response.Empreendimento = &EmpreendimentoResponse{
			ID:              imovel.Empreendimento.ID,
			Titulo:          imovel.Empreendimento.Titulo,
			Descricao:       imovel.Empreendimento.Descricao,
			DataEntrega:     imovel.Empreendimento.DataEntrega,
			EtapaLancamento: imovel.Empreendimento.EtapaLancamento,
			Finalidade:      imovel.Empreendimento.Finalidade,
			Tipo:            imovel.Empreendimento.Tipo,
			Status:          imovel.Empreendimento.Status,
			Localizacao:     imovel.Empreendimento.Localizacao,
			CreatedAt:       imovel.Empreendimento.CreatedAt,
			UpdatedAt:       imovel.Empreendimento.UpdatedAt,
		}
	}

	if imovel.Planta != nil {
		response.Planta = &PlantaResponse{
			ID:        imovel.Planta.ID,
			Nome:      imovel.Planta.Nome,
			Metragem:  imovel.Planta.Metragem,
			CreatedAt: imovel.Planta.CreatedAt,
			UpdatedAt: imovel.Planta.UpdatedAt,
		}
	}

	if imovel.CorretorPrincipal != nil {
		response.CorretorPrincipal = &CorretorPrincipalResponse{
			ID:             imovel.CorretorPrincipal.ID,
			Nome:           imovel.CorretorPrincipal.Nome,
			Email:          imovel.CorretorPrincipal.Email,
			Whatsapp:       imovel.CorretorPrincipal.Whatsapp,
			Idiomas:        imovel.CorretorPrincipal.Idiomas,
			BairrosAtuacao: imovel.CorretorPrincipal.BairrosAtuacao,
		}

		// Map Foto if present
		if imovel.CorretorPrincipal.Foto != nil {
			response.CorretorPrincipal.Foto = &AnexoResponse{
				ID:            imovel.CorretorPrincipal.Foto.ID,
				Nome:          imovel.CorretorPrincipal.Foto.Nome,
				Path:          imovel.CorretorPrincipal.Foto.Path,
				Tamanho:       imovel.CorretorPrincipal.Foto.Tamanho,
				Tipo:          imovel.CorretorPrincipal.Foto.Tipo,
				URL:           imovel.CorretorPrincipal.Foto.URL,
				CanPublish:    imovel.CorretorPrincipal.Foto.CanPublish,
				Image:         imovel.CorretorPrincipal.Foto.Image,
				Video:         imovel.CorretorPrincipal.Foto.Video,
				IsExternalURL: imovel.CorretorPrincipal.Foto.IsExternalURL,
				CreatedAt:     imovel.CorretorPrincipal.Foto.CreatedAt,
				UpdatedAt:     imovel.CorretorPrincipal.Foto.UpdatedAt,
			}
		}

		// Map Organizacao if present
		if imovel.CorretorPrincipal.Organizacao != nil {
			response.CorretorPrincipal.Organizacao = &OrganizacaoResponse{
				ID:     imovel.CorretorPrincipal.Organizacao.ID,
				Nome:   imovel.CorretorPrincipal.Organizacao.Nome,
				Perfil: imovel.CorretorPrincipal.Organizacao.Perfil,
			}
		}
	}

	if imovel.Pacote != nil {
		response.Pacote = &PacoteResponse{
			ID:         imovel.Pacote.ID,
			Titulo:     imovel.Pacote.Titulo,
			Descricao:  imovel.Pacote.Descricao,
			Exclusivo:  imovel.Pacote.Exclusivo,
			EmDestaque: imovel.Pacote.EmDestaque,
			CreatedAt:  imovel.Pacote.CreatedAt,
			UpdatedAt:  imovel.Pacote.UpdatedAt,
		}
	}

	if imovel.PrecoVenda != nil {
		response.PrecoVenda = &PrecoVendaResponse{
			ID:                          imovel.PrecoVenda.ID,
			Preco:                       imovel.PrecoVenda.Preco,
			AceitaFinanciamentoBancario: imovel.PrecoVenda.AceitaFinanciamentoBancario,
			AceitaFinanciamentoDireto:   imovel.PrecoVenda.AceitaFinanciamentoDireto,
			AceitaPermuta:               imovel.PrecoVenda.AceitaPermuta,
			AceitaCartaDeCredito:        imovel.PrecoVenda.AceitaCartaDeCredito,
			AceitaFGTS:                  imovel.PrecoVenda.AceitaFGTS,
			Ativo:                       imovel.PrecoVenda.Ativo,
			PacoteTitulo:                imovel.PrecoVenda.PacoteTitulo,
			PacoteDescricao:             imovel.PrecoVenda.PacoteDescricao,
			PacoteExclusivo:             imovel.PrecoVenda.PacoteExclusivo,
			PacoteEmDestaque:            imovel.PrecoVenda.PacoteEmDestaque,
			CreatedAt:                   imovel.PrecoVenda.CreatedAt,
			UpdatedAt:                   imovel.PrecoVenda.UpdatedAt,
		}
	}

	if imovel.PrecoAluguel != nil {
		response.PrecoAluguel = &PrecoAluguelResponse{
			ID:           imovel.PrecoAluguel.ID,
			Preco:        imovel.PrecoAluguel.Preco,
			AceitaFiador: imovel.PrecoAluguel.AceitaFiador,
			Ativo:        imovel.PrecoAluguel.Ativo,
			CreatedAt:    imovel.PrecoAluguel.CreatedAt,
			UpdatedAt:    imovel.PrecoAluguel.UpdatedAt,
		}
	}

	// Map anexos
	if len(imovel.Anexos) > 0 {
		response.Anexos = make([]AnexoResponse, len(imovel.Anexos))
		for i, anexo := range imovel.Anexos {
			response.Anexos[i] = AnexoResponse{
				ID:            anexo.ID,
				Nome:          anexo.Nome,
				Path:          anexo.Path,
				Tamanho:       anexo.Tamanho,
				Tipo:          anexo.Tipo,
				URL:           anexo.URL,
				CanPublish:    anexo.CanPublish,
				Image:         anexo.Image,
				Video:         anexo.Video,
				IsExternalURL: anexo.IsExternalURL,
				CreatedAt:     anexo.CreatedAt,
				UpdatedAt:     anexo.UpdatedAt,
			}
		}
	}

	return response
}

// Relationship Management Methods

// AddAnexo adds an attachment to a property
func (s *service) AddAnexo(ctx context.Context, imovelID uint, anexo *Anexo) error {
	if imovelID == 0 {
		return errors.New("invalid property ID")
	}

	imovel, err := s.repo.FindByID(ctx, imovelID)
	if err != nil {
		return fmt.Errorf("failed to find property: %w", err)
	}

	if imovel == nil {
		return fmt.Errorf("property not found")
	}

	if err := s.repo.AddAnexo(ctx, imovelID, anexo); err != nil {
		return fmt.Errorf("failed to add attachment: %w", err)
	}

	return nil
}

// RemoveAnexo removes an attachment from a property
func (s *service) RemoveAnexo(ctx context.Context, imovelID, anexoID uint) error {
	if imovelID == 0 || anexoID == 0 {
		return errors.New("invalid property or attachment ID")
	}

	if err := s.repo.RemoveAnexo(ctx, imovelID, anexoID); err != nil {
		return fmt.Errorf("failed to remove attachment: %w", err)
	}

	return nil
}

// GetAnexos retrieves all attachments for a property
func (s *service) GetAnexos(ctx context.Context, imovelID uint) ([]AnexoResponse, error) {
	if imovelID == 0 {
		return nil, errors.New("invalid property ID")
	}

	anexos, err := s.repo.GetAnexos(ctx, imovelID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve attachments: %w", err)
	}

	responses := make([]AnexoResponse, len(anexos))
	for i, anexo := range anexos {
		responses[i] = AnexoResponse{
			ID:            anexo.ID,
			Nome:          anexo.Nome,
			Path:          anexo.Path,
			Tamanho:       anexo.Tamanho,
			Tipo:          anexo.Tipo,
			URL:           anexo.URL,
			CanPublish:    anexo.CanPublish,
			Image:         anexo.Image,
			Video:         anexo.Video,
			IsExternalURL: anexo.IsExternalURL,
			CreatedAt:     anexo.CreatedAt,
			UpdatedAt:     anexo.UpdatedAt,
		}
	}

	return responses, nil
}

// AttachEndereco attaches an address to a property
func (s *service) AttachEndereco(ctx context.Context, imovelID, enderecoID uint) error {
	if imovelID == 0 || enderecoID == 0 {
		return errors.New("invalid property or address ID")
	}

	imovel, err := s.repo.FindByID(ctx, imovelID)
	if err != nil {
		return fmt.Errorf("failed to find property: %w", err)
	}

	if imovel == nil {
		return fmt.Errorf("property not found")
	}

	if err := s.repo.UpdateEndereco(ctx, imovelID, enderecoID); err != nil {
		return fmt.Errorf("failed to attach address: %w", err)
	}

	return nil
}

// CreateEndereco creates a new address
func (s *service) CreateEndereco(ctx context.Context, endereco *Endereco) error {
	return s.repo.CreateEndereco(ctx, endereco)
}

// AttachEmpreendimento attaches an enterprise to a property
func (s *service) AttachEmpreendimento(ctx context.Context, imovelID, empreendimentoID uint) error {
	if imovelID == 0 || empreendimentoID == 0 {
		return errors.New("invalid property or enterprise ID")
	}

	imovel, err := s.repo.FindByID(ctx, imovelID)
	if err != nil {
		return fmt.Errorf("failed to find property: %w", err)
	}

	if imovel == nil {
		return fmt.Errorf("property not found")
	}

	if err := s.repo.UpdateEmpreendimento(ctx, imovelID, empreendimentoID); err != nil {
		return fmt.Errorf("failed to attach enterprise: %w", err)
	}

	return nil
}

// AttachPlanta attaches a floor plan to a property
func (s *service) AttachPlanta(ctx context.Context, imovelID, plantaID uint) error {
	if imovelID == 0 || plantaID == 0 {
		return errors.New("invalid property or floor plan ID")
	}

	imovel, err := s.repo.FindByID(ctx, imovelID)
	if err != nil {
		return fmt.Errorf("failed to find property: %w", err)
	}

	if imovel == nil {
		return fmt.Errorf("property not found")
	}

	if err := s.repo.UpdatePlanta(ctx, imovelID, plantaID); err != nil {
		return fmt.Errorf("failed to attach floor plan: %w", err)
	}

	return nil
}

// AttachPacote attaches a package to a property
func (s *service) AttachPacote(ctx context.Context, imovelID, pacoteID uint) error {
	if imovelID == 0 || pacoteID == 0 {
		return errors.New("invalid property or package ID")
	}

	imovel, err := s.repo.FindByID(ctx, imovelID)
	if err != nil {
		return fmt.Errorf("failed to find property: %w", err)
	}

	if imovel == nil {
		return fmt.Errorf("property not found")
	}

	if err := s.repo.UpdatePacote(ctx, imovelID, pacoteID); err != nil {
		return fmt.Errorf("failed to attach package: %w", err)
	}

	return nil
}

// AttachOrganizacao attaches an organization to a property
func (s *service) AttachOrganizacao(ctx context.Context, imovelID, organizacaoID uint) error {
	if imovelID == 0 || organizacaoID == 0 {
		return errors.New("invalid property or organization ID")
	}

	imovel, err := s.repo.FindByID(ctx, imovelID)
	if err != nil {
		return fmt.Errorf("failed to find property: %w", err)
	}

	if imovel == nil {
		return fmt.Errorf("property not found")
	}

	if err := s.repo.UpdateCorretorPrincipal(ctx, imovelID, organizacaoID); err != nil {
		return fmt.Errorf("failed to attach organization: %w", err)
	}

	return nil
}

// AttachPrecoVenda attaches a selling price to a property
func (s *service) AttachPrecoVenda(ctx context.Context, imovelID, precoVendaID uint) error {
	if imovelID == 0 || precoVendaID == 0 {
		return errors.New("invalid property or price ID")
	}

	imovel, err := s.repo.FindByID(ctx, imovelID)
	if err != nil {
		return fmt.Errorf("failed to find property: %w", err)
	}

	if imovel == nil {
		return fmt.Errorf("property not found")
	}

	if err := s.repo.UpdatePrecoVenda(ctx, imovelID, precoVendaID); err != nil {
		return fmt.Errorf("failed to attach selling price: %w", err)
	}

	return nil
}

// AttachPrecoAluguel attaches a rental price to a property
func (s *service) AttachPrecoAluguel(ctx context.Context, imovelID, precoAluguelID uint) error {
	if imovelID == 0 || precoAluguelID == 0 {
		return errors.New("invalid property or price ID")
	}

	imovel, err := s.repo.FindByID(ctx, imovelID)
	if err != nil {
		return fmt.Errorf("failed to find property: %w", err)
	}

	if imovel == nil {
		return fmt.Errorf("property not found")
	}

	if err := s.repo.UpdatePrecoAluguel(ctx, imovelID, precoAluguelID); err != nil {
		return fmt.Errorf("failed to attach rental price: %w", err)
	}

	return nil
}

// AddCaracteristicas adds characteristics to a property
func (s *service) AddCaracteristicas(ctx context.Context, imovelID uint, caracteristicaIDs []uint) error {
	if imovelID == 0 {
		return errors.New("invalid property ID")
	}

	if len(caracteristicaIDs) == 0 {
		return nil
	}

	imovel, err := s.repo.FindByID(ctx, imovelID)
	if err != nil {
		return fmt.Errorf("failed to find property: %w", err)
	}

	if imovel == nil {
		return fmt.Errorf("property not found")
	}

	if err := s.repo.AddCaracteristicas(ctx, imovelID, caracteristicaIDs); err != nil {
		return fmt.Errorf("failed to add characteristics: %w", err)
	}

	return nil
}

// RemoveCaracteristicas removes characteristics from a property
func (s *service) RemoveCaracteristicas(ctx context.Context, imovelID uint, caracteristicaIDs []uint) error {
	if imovelID == 0 {
		return errors.New("invalid property ID")
	}

	if len(caracteristicaIDs) == 0 {
		return nil
	}

	imovel, err := s.repo.FindByID(ctx, imovelID)
	if err != nil {
		return fmt.Errorf("failed to find property: %w", err)
	}

	if imovel == nil {
		return fmt.Errorf("property not found")
	}

	if err := s.repo.RemoveCaracteristicas(ctx, imovelID, caracteristicaIDs); err != nil {
		return fmt.Errorf("failed to remove characteristics: %w", err)
	}

	return nil
}

// GetCaracteristicas retrieves all characteristics for a property
func (s *service) GetCaracteristicas(ctx context.Context, imovelID uint) ([]CaracteristicaResponse, error) {
	if imovelID == 0 {
		return nil, errors.New("invalid property ID")
	}

	caracteristicas, err := s.repo.GetCaracteristicas(ctx, imovelID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve characteristics: %w", err)
	}

	responses := make([]CaracteristicaResponse, len(caracteristicas))
	for i, caract := range caracteristicas {
		responses[i] = CaracteristicaResponse{
			ID:            caract.ID,
			Nome:          caract.Nome,
			CategoriaID:   caract.CategoriaID,
			CategoriaNome: caract.CategoriaNome,
			CreatedAt:     caract.CreatedAt,
			UpdatedAt:     caract.UpdatedAt,
		}
	}

	return responses, nil
}

// ReplaceCaracteristicas replaces all characteristics for a property
func (s *service) ReplaceCaracteristicas(ctx context.Context, imovelID uint, caracteristicaIDs []uint) error {
	if imovelID == 0 {
		return errors.New("invalid property ID")
	}

	imovel, err := s.repo.FindByID(ctx, imovelID)
	if err != nil {
		return fmt.Errorf("failed to find property: %w", err)
	}

	if imovel == nil {
		return fmt.Errorf("property not found")
	}

	// Remove all existing characteristics
	if err := s.repo.RemoveAllCaracteristicas(ctx, imovelID); err != nil {
		return fmt.Errorf("failed to remove existing characteristics: %w", err)
	}

	// Add new characteristics
	if len(caracteristicaIDs) > 0 {
		if err := s.repo.AddCaracteristicas(ctx, imovelID, caracteristicaIDs); err != nil {
			return fmt.Errorf("failed to add characteristics: %w", err)
		}
	}

	return nil
}
