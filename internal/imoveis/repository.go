package imoveis

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Repository defines the interface for property data access
type Repository interface {
	// Create
	Create(ctx context.Context, imovel *Imovel) error

	// Read
	FindByID(ctx context.Context, id uint) (*Imovel, error)
	FindByCodigo(ctx context.Context, codigo string) (*Imovel, error)
	FindByIdIntegracao(ctx context.Context, idIntegracao string) (*Imovel, error)

	// Update
	Update(ctx context.Context, imovel *Imovel) error

	// Delete
	Delete(ctx context.Context, id uint) error
	HardDelete(ctx context.Context, id uint) error

	// List & Filter
	List(ctx context.Context, query *ImovelListQuery) (*ImovelListResponse, error)
	ListByEmpreendimento(ctx context.Context, empreendimentoID uint, page, limit int) ([]Imovel, int64, error)
	ListByCorretorPrincipal(ctx context.Context, corretorPrincipalID uint, page, limit int) ([]Imovel, int64, error)

	// Bulk Operations
	CreateBatch(ctx context.Context, imoveis []Imovel) error
	UpdateBatch(ctx context.Context, imoveis []Imovel) error

	// Count
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
	CountByEmpreendimento(ctx context.Context, empreendimentoID uint) (int64, error)

	// Exists
	ExistsByCodigo(ctx context.Context, codigo string) (bool, error)
	ExistsByIdIntegracao(ctx context.Context, idIntegracao string) (bool, error)

	// Relationships - Anexos
	AddAnexo(ctx context.Context, imovelID uint, anexo *Anexo) error
	RemoveAnexo(ctx context.Context, imovelID, anexoID uint) error
	GetAnexos(ctx context.Context, imovelID uint) ([]Anexo, error)

	// Relationships - Single associations
	UpdateEndereco(ctx context.Context, imovelID, enderecoID uint) error
	UpdateEmpreendimento(ctx context.Context, imovelID, empreendimentoID uint) error
	UpdatePlanta(ctx context.Context, imovelID, plantaID uint) error
	UpdatePacote(ctx context.Context, imovelID, pacoteID uint) error
	UpdateCorretorPrincipal(ctx context.Context, imovelID, corretorPrincipalID uint) error
	UpdatePrecoVenda(ctx context.Context, imovelID, precoVendaID uint) error
	UpdatePrecoAluguel(ctx context.Context, imovelID, precoAluguelID uint) error

	// Endereco management
	CreateEndereco(ctx context.Context, endereco *Endereco) error

	// Relationships - Caracteristicas
	AddCaracteristicas(ctx context.Context, imovelID uint, caracteristicaIDs []uint) error
	RemoveCaracteristicas(ctx context.Context, imovelID uint, caracteristicaIDs []uint) error
	GetCaracteristicas(ctx context.Context, imovelID uint) ([]Caracteristica, error)
	RemoveAllCaracteristicas(ctx context.Context, imovelID uint) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new property repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new property
func (r *repository) Create(ctx context.Context, imovel *Imovel) error {
	if err := r.db.WithContext(ctx).Create(imovel).Error; err != nil {
		return err
	}
	return nil
}

// FindByID retrieves a property by ID with all relations
func (r *repository) FindByID(ctx context.Context, id uint) (*Imovel, error) {
	var imovel Imovel
	if err := r.db.WithContext(ctx).
		Preload("Endereco").
		Preload("Empreendimento", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Endereco").Preload("Torres").Preload("Plantas").Preload("Caracteristicas").Preload("Anexos")
		}).
		Preload("Planta", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Anexos")
		}).
		Preload("CorretorPrincipal").
		Preload("CorretorPrincipal.Organizacao").
		Preload("CorretorPrincipal.Foto").
		Preload("Pacote").
		Preload("PrecoVenda").
		Preload("PrecoAluguel").
		Preload("Anexos").
		Where("id = ?", id).
		First(&imovel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &imovel, nil
}

// FindByCodigo retrieves a property by codigo
func (r *repository) FindByCodigo(ctx context.Context, codigo string) (*Imovel, error) {
	var imovel Imovel
	if err := r.db.WithContext(ctx).
		Preload("Endereco").
		Preload("Empreendimento", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Endereco").Preload("Torres").Preload("Plantas").Preload("Caracteristicas").Preload("Anexos")
		}).
		Preload("Planta", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Anexos")
		}).
		Preload("CorretorPrincipal").
		Preload("CorretorPrincipal.Organizacao").
		Preload("CorretorPrincipal.Foto").
		Preload("Pacote").
		Preload("PrecoVenda").
		Preload("PrecoAluguel").
		Preload("Anexos").
		Where("codigo = ?", codigo).
		First(&imovel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &imovel, nil
}

// FindByIdIntegracao retrieves a property by integration ID
func (r *repository) FindByIdIntegracao(ctx context.Context, idIntegracao string) (*Imovel, error) {
	var imovel Imovel
	if err := r.db.WithContext(ctx).
		Preload("Endereco").
		Preload("Empreendimento", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Endereco").Preload("Torres").Preload("Plantas").Preload("Caracteristicas").Preload("Anexos")
		}).
		Preload("Planta", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Anexos")
		}).
		Preload("CorretorPrincipal").
		Preload("CorretorPrincipal.Organizacao").
		Preload("CorretorPrincipal.Foto").
		Preload("Pacote").
		Preload("PrecoVenda").
		Preload("PrecoAluguel").
		Preload("Anexos").
		Where("id_integracao = ?", idIntegracao).
		First(&imovel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &imovel, nil
}

// Update updates a property
func (r *repository) Update(ctx context.Context, imovel *Imovel) error {
	// Omit associations to prevent GORM from trying to update them
	// Only update the imovel table fields, not related entities
	if err := r.db.WithContext(ctx).Model(imovel).
		Omit("Endereco", "Empreendimento", "Planta", "CorretorPrincipal", "Pacote", "PrecoVenda", "PrecoAluguel", "Anexos").
		Updates(imovel).Error; err != nil {
		return err
	}
	return nil
}

// Delete soft deletes a property
func (r *repository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&Imovel{}, id).Error; err != nil {
		return err
	}
	return nil
}

// HardDelete permanently deletes a property
func (r *repository) HardDelete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Unscoped().Delete(&Imovel{}, id).Error; err != nil {
		return err
	}
	return nil
}

// List retrieves properties with filtering and pagination
func (r *repository) List(ctx context.Context, query *ImovelListQuery) (*ImovelListResponse, error) {
	var imoveis []Imovel
	var total int64

	db := r.db.WithContext(ctx)

	// Apply filters
	if query.Codigo != "" {
		db = db.Where("codigo ILIKE ?", "%"+query.Codigo+"%")
	}
	if query.Tipo != "" {
		db = db.Where("tipo = ?", query.Tipo)
	}
	if query.Objetivo != "" {
		db = db.Where("objetivo = ?", query.Objetivo)
	}
	if query.Finalidade != "" {
		db = db.Where("finalidade = ?", query.Finalidade)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Published != nil {
		db = db.Where("published = ?", *query.Published)
	}
	if query.MinPreco > 0 {
		db = db.Joins("LEFT JOIN preco_vendas ON preco_vendas.id = imoveis.preco_venda_id").
			Where("preco_vendas.preco >= ?", query.MinPreco)
	}
	if query.MaxPreco > 0 {
		db = db.Joins("LEFT JOIN preco_vendas ON preco_vendas.id = imoveis.preco_venda_id").
			Where("preco_vendas.preco <= ?", query.MaxPreco)
	}
	if query.MinMetragem > 0 {
		db = db.Where("metragem >= ?", query.MinMetragem)
	}
	if query.MaxMetragem > 0 {
		db = db.Where("metragem <= ?", query.MaxMetragem)
	}
	if query.Rua != "" {
		db = db.Joins("INNER JOIN enderecos ON enderecos.id = imoveis.endereco_id").
			Where("enderecos.rua ILIKE ?", "%"+query.Rua+"%")
	}
	if query.Cidade != "" {
		db = db.Joins("INNER JOIN enderecos ON enderecos.id = imoveis.endereco_id").
			Where("enderecos.cidade ILIKE ?", "%"+query.Cidade+"%")
	}
	if query.Bairro != "" {
		db = db.Joins("INNER JOIN enderecos ON enderecos.id = imoveis.endereco_id").
			Where("enderecos.bairro ILIKE ?", "%"+query.Bairro+"%")
	}
	if query.NumQuartos > 0 {
		db = db.Where("num_quartos >= ?", query.NumQuartos)
	}
	if query.NumBanheiros > 0 {
		db = db.Where("num_banheiros >= ?", query.NumBanheiros)
	}
	if query.NumGaragens > 0 {
		db = db.Where("num_vagas >= ?", query.NumGaragens)
	}
	if query.EmpreendimentoID > 0 {
		db = db.Where("empreendimento_id = ?", query.EmpreendimentoID)
	}

	// Count total
	if err := db.Model(&Imovel{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	sortField := "created_at"
	if query.Sort != "" {
		sortField = query.Sort
	}
	order := "DESC"
	if query.Order == "asc" {
		order = "ASC"
	}
	db = db.Order(sortField + " " + order)

	// Apply pagination
	offset := (query.Page - 1) * query.Limit
	if err := db.Preload("Endereco").
		Preload("Empreendimento", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Endereco")
		}).
		Preload("Planta", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Anexos")
		}).
		Preload("CorretorPrincipal").
		Preload("CorretorPrincipal.Organizacao").
		Preload("CorretorPrincipal.Foto").
		Preload("Pacote").
		Preload("PrecoVenda").
		Preload("PrecoAluguel").
		Preload("Anexos").
		Offset(offset).
		Limit(query.Limit).
		Find(&imoveis).Error; err != nil {
		return nil, err
	}

	// Build response
	pages := (total + int64(query.Limit) - 1) / int64(query.Limit)
	results := make([]ImovelResponse, len(imoveis))
	for i, imovel := range imoveis {
		results[i] = r.mapToResponse(&imovel)
	}

	return &ImovelListResponse{
		Total:   total,
		Page:    query.Page,
		Limit:   query.Limit,
		Pages:   pages,
		HasNext: int64(query.Page) < pages,
		HasPrev: query.Page > 1,
		Results: results,
	}, nil
}

// ListByEmpreendimento retrieves properties by enterprise
func (r *repository) ListByEmpreendimento(ctx context.Context, empreendimentoID uint, page, limit int) ([]Imovel, int64, error) {
	var imoveis []Imovel
	var total int64

	db := r.db.WithContext(ctx).Where("empreendimento_id = ?", empreendimentoID)

	if err := db.Model(&Imovel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := db.Preload("Endereco").
		Preload("Empreendimento").
		Preload("Planta").
		Preload("CorretorPrincipal").
		Preload("CorretorPrincipal.Organizacao").
		Preload("CorretorPrincipal.Foto").
		Preload("Pacote").
		Preload("PrecoVenda").
		Preload("PrecoAluguel").
		Preload("Anexos").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&imoveis).Error; err != nil {
		return nil, 0, err
	}

	return imoveis, total, nil
}

// ListByCorretorPrincipal retrieves properties by real estate agent
func (r *repository) ListByCorretorPrincipal(ctx context.Context, corretorPrincipalID uint, page, limit int) ([]Imovel, int64, error) {
	var imoveis []Imovel
	var total int64

	db := r.db.WithContext(ctx).Where("corretor_principal_id = ?", corretorPrincipalID)

	if err := db.Model(&Imovel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := db.Preload("Endereco").
		Preload("Empreendimento").
		Preload("Planta").
		Preload("CorretorPrincipal").
		Preload("CorretorPrincipal.Organizacao").
		Preload("Pacote").
		Preload("PrecoVenda").
		Preload("PrecoAluguel").
		Preload("Anexos").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&imoveis).Error; err != nil {
		return nil, 0, err
	}

	return imoveis, total, nil
}

// CreateBatch creates multiple properties
func (r *repository) CreateBatch(ctx context.Context, imoveis []Imovel) error {
	if err := r.db.WithContext(ctx).CreateInBatches(imoveis, 100).Error; err != nil {
		return err
	}
	return nil
}

// UpdateBatch updates multiple properties
func (r *repository) UpdateBatch(ctx context.Context, imoveis []Imovel) error {
	if err := r.db.WithContext(ctx).Save(imoveis).Error; err != nil {
		return err
	}
	return nil
}

// Count returns total number of properties
func (r *repository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&Imovel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByStatus returns count of properties by status
func (r *repository) CountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&Imovel{}).
		Where("status = ?", status).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByEmpreendimento returns count of properties by enterprise
func (r *repository) CountByEmpreendimento(ctx context.Context, empreendimentoID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&Imovel{}).
		Where("empreendimento_id = ?", empreendimentoID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ExistsByCodigo checks if a property exists by codigo
func (r *repository) ExistsByCodigo(ctx context.Context, codigo string) (bool, error) {
	var exists bool
	if err := r.db.WithContext(ctx).
		Model(&Imovel{}).
		Select("count(*) > 0").
		Where("codigo = ?", codigo).
		Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

// ExistsByIdIntegracao checks if a property exists by integration ID
func (r *repository) ExistsByIdIntegracao(ctx context.Context, idIntegracao string) (bool, error) {
	var exists bool
	if err := r.db.WithContext(ctx).
		Model(&Imovel{}).
		Select("count(*) > 0").
		Where("id_integracao = ?", idIntegracao).
		Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

// AddAnexo adds an attachment to a property
func (r *repository) AddAnexo(ctx context.Context, imovelID uint, anexo *Anexo) error {
	imovelIDPtr := &imovelID
	anexo.ImovelID = imovelIDPtr

	// Build list of fields to omit based on nil values
	var omitFields []string
	if anexo.EmpreendimentoID == nil {
		omitFields = append(omitFields, "EmpreendimentoID")
	}
	if anexo.PlantaID == nil {
		omitFields = append(omitFields, "PlantaID")
	}

	// Create anexo, omitting zero-value foreign keys to avoid constraint violations
	db := r.db.WithContext(ctx)
	if len(omitFields) > 0 {
		db = db.Omit(omitFields...)
	}

	if err := db.Create(anexo).Error; err != nil {
		return err
	}
	return nil
}

// RemoveAnexo removes an attachment from a property
func (r *repository) RemoveAnexo(ctx context.Context, imovelID, anexoID uint) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND imovel_id = ?", anexoID, imovelID).Delete(&Anexo{}).Error; err != nil {
		return err
	}
	return nil
}

// GetAnexos retrieves all attachments for a property
func (r *repository) GetAnexos(ctx context.Context, imovelID uint) ([]Anexo, error) {
	var anexos []Anexo
	if err := r.db.WithContext(ctx).
		Where("imovel_id = ?", imovelID).
		Order("created_at DESC").
		Find(&anexos).Error; err != nil {
		return nil, err
	}
	return anexos, nil
}

// UpdateEndereco updates the address of a property
func (r *repository) UpdateEndereco(ctx context.Context, imovelID, enderecoID uint) error {
	if err := r.db.WithContext(ctx).Model(&Imovel{}).
		Where("id = ?", imovelID).
		Update("endereco_id", enderecoID).Error; err != nil {
		return err
	}
	return nil
}

// UpdateEmpreendimento updates the enterprise of a property
func (r *repository) UpdateEmpreendimento(ctx context.Context, imovelID, empreendimentoID uint) error {
	if err := r.db.WithContext(ctx).Model(&Imovel{}).
		Where("id = ?", imovelID).
		Update("empreendimento_id", empreendimentoID).Error; err != nil {
		return err
	}
	return nil
}

// UpdatePlanta updates the floor plan of a property
func (r *repository) UpdatePlanta(ctx context.Context, imovelID, plantaID uint) error {
	if err := r.db.WithContext(ctx).Model(&Imovel{}).
		Where("id = ?", imovelID).
		Update("planta_id", plantaID).Error; err != nil {
		return err
	}
	return nil
}

// UpdatePacote updates the package of a property
func (r *repository) UpdatePacote(ctx context.Context, imovelID, pacoteID uint) error {
	if err := r.db.WithContext(ctx).Model(&Imovel{}).
		Where("id = ?", imovelID).
		Update("pacote_id", pacoteID).Error; err != nil {
		return err
	}
	return nil
}

// UpdateCorretorPrincipal updates the real estate agent of a property
func (r *repository) UpdateCorretorPrincipal(ctx context.Context, imovelID, corretorPrincipalID uint) error {
	if err := r.db.WithContext(ctx).Model(&Imovel{}).
		Where("id = ?", imovelID).
		Update("corretor_principal_id", corretorPrincipalID).Error; err != nil {
		return err
	}
	return nil
}

// UpdatePrecoVenda updates the selling price of a property
func (r *repository) UpdatePrecoVenda(ctx context.Context, imovelID, precoVendaID uint) error {
	if err := r.db.WithContext(ctx).Model(&Imovel{}).
		Where("id = ?", imovelID).
		Update("preco_venda_id", precoVendaID).Error; err != nil {
		return err
	}
	return nil
}

// UpdatePrecoAluguel updates the rental price of a property
func (r *repository) UpdatePrecoAluguel(ctx context.Context, imovelID, precoAluguelID uint) error {
	if err := r.db.WithContext(ctx).Model(&Imovel{}).
		Where("id = ?", imovelID).
		Update("preco_aluguel_id", precoAluguelID).Error; err != nil {
		return err
	}
	return nil
}

// AddCaracteristicas adds characteristics to a property
func (r *repository) AddCaracteristicas(ctx context.Context, imovelID uint, caracteristicaIDs []uint) error {
	if len(caracteristicaIDs) == 0 {
		return nil
	}

	imovel := &Imovel{ID: imovelID}
	caracteristicas := make([]Caracteristica, len(caracteristicaIDs))
	for i, id := range caracteristicaIDs {
		caracteristicas[i] = Caracteristica{ID: id}
	}

	if err := r.db.WithContext(ctx).Model(imovel).Association("Caracteristicas").Append(caracteristicas); err != nil {
		return err
	}
	return nil
}

// RemoveCaracteristicas removes characteristics from a property
func (r *repository) RemoveCaracteristicas(ctx context.Context, imovelID uint, caracteristicaIDs []uint) error {
	if len(caracteristicaIDs) == 0 {
		return nil
	}

	imovel := &Imovel{ID: imovelID}
	if err := r.db.WithContext(ctx).Model(imovel).Association("Caracteristicas").Delete(caracteristicaIDs); err != nil {
		return err
	}
	return nil
}

// GetCaracteristicas retrieves all characteristics for a property
func (r *repository) GetCaracteristicas(ctx context.Context, imovelID uint) ([]Caracteristica, error) {
	var caracteristicas []Caracteristica
	if err := r.db.WithContext(ctx).
		Model(&Imovel{ID: imovelID}).
		Association("Caracteristicas").
		Find(&caracteristicas); err != nil {
		return nil, err
	}
	return caracteristicas, nil
}

// RemoveAllCaracteristicas removes all characteristics from a property
func (r *repository) RemoveAllCaracteristicas(ctx context.Context, imovelID uint) error {
	if err := r.db.WithContext(ctx).
		Model(&Imovel{ID: imovelID}).
		Association("Caracteristicas").
		Clear(); err != nil {
		return err
	}
	return nil
}

// mapToResponse converts Imovel model to response DTO
func (r *repository) mapToResponse(imovel *Imovel) ImovelResponse {
	response := ImovelResponse{
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

// CreateEndereco creates a new address
func (r *repository) CreateEndereco(ctx context.Context, endereco *Endereco) error {
	return r.db.WithContext(ctx).Create(endereco).Error
}
