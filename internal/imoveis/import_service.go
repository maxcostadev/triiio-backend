package imoveis

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/vahiiiid/go-rest-api-boilerplate/internal/config"
)

// ImportService defines the interface for importing properties from external API
type ImportService interface {
	ImportPublishedProperties(ctx context.Context) error
	ImportPropertyDetails(ctx context.Context, externalID uint) (*ExternalDetailedImovel, error)
}

type importService struct {
	service           Service
	httpClient        *http.Client
	baseURL           string
	apiKey            string
	integrationSource string
}

// NewImportService creates a new import service
func NewImportService(service Service, extCfg *config.ExternalAPIConfig) ImportService {
	timeout := time.Duration(extCfg.TimeoutSeconds) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &importService{
		service:           service,
		httpClient:        &http.Client{Timeout: timeout},
		baseURL:           extCfg.BaseURL,
		apiKey:            extCfg.APIKey,
		integrationSource: extCfg.IntegrationSource,
	}
}

// ImportPublishedProperties imports all published properties from external API
// Uses upsert logic: creates new properties or updates existing ones
func (is *importService) ImportPublishedProperties(ctx context.Context) error {
	// Fetch list of published properties
	listURL := fmt.Sprintf("%s/api/properties/published", is.baseURL)

	properties, err := is.fetchPublishedList(ctx, listURL)
	if err != nil {
		return fmt.Errorf("failed to fetch published properties: %w", err)
	}

	if len(properties) == 0 {
		return fmt.Errorf("no properties found in external API")
	}

	// Process each property
	var successCount, errorCount, updateCount int
	for _, extImovel := range properties {
		// Fetch detailed info for this property (includes empreendimento and torres)
		log.Printf("####PROPERTIER %v", extImovel.ID)
		detailedImovel, err := is.ImportPropertyDetails(ctx, extImovel.ID)
		if err != nil {
			fmt.Printf("Warning: Failed to fetch details for property %d: %v\n", extImovel.ID, err)
			errorCount++
			continue
		}

		idIntegracao := fmt.Sprintf("%d", detailedImovel.ID)

		// Check if property already exists by IdIntegracao
		existingImovel, err := is.service.GetImovelByIdIntegracao(ctx, idIntegracao)
		if err == nil && existingImovel != nil {
			// Property exists - update it and its relationships
			fmt.Printf("Property %s already exists (ID: %d), updating...\n", detailedImovel.Codigo, existingImovel.ID)
			if _, err := is.upsertImovelAndRelationships(ctx, existingImovel.ID, detailedImovel, true); err != nil {
				fmt.Printf("Warning: Failed to update property %s: %v\n", detailedImovel.Codigo, err)
				errorCount++
				continue
			}
			updateCount++
		} else {
			// Property doesn't exist - create it and its relationships
			imovelResp, err := is.upsertImovelAndRelationships(ctx, 0, detailedImovel, false)
			if err != nil {
				fmt.Printf("Warning: Failed to create property %s: %v\n", detailedImovel.Codigo, err)
				errorCount++
				continue
			}

			fmt.Printf("Successfully created property: %s (ID: %d)\n", detailedImovel.Codigo, imovelResp.ID)
			successCount++
		}
	}

	return fmt.Errorf("import completed: %d created, %d updated, %d failed", successCount, updateCount, errorCount)
}

// ImportPropertyDetails fetches detailed property information including empreendimento
func (is *importService) ImportPropertyDetails(ctx context.Context, externalID uint) (*ExternalDetailedImovel, error) {
	detailURL := fmt.Sprintf("%s/api/properties/published/%d", is.baseURL, externalID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, detailURL, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	is.setHeaders(req)

	resp, err := is.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch property details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Results ExternalDetailedImovel `json:"results"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result.Results, nil
}

// upsertImovelAndRelationships creates or updates a property and all its relationships
// isUpdate=true means we're updating an existing property, false means creating new
func (is *importService) upsertImovelAndRelationships(ctx context.Context, imovelID uint, ext *ExternalDetailedImovel, isUpdate bool) (*ImovelResponse, error) {
	var imovelResp *ImovelResponse
	var err error

	// Always upsert relationships first (works for both create and update)
	var empreendimentoID uint
	if ext.Empreendimento != nil {
		empID, err := is.upsertEmpreendimento(ctx, ext.Empreendimento)
		if err != nil {
			fmt.Printf("Warning: Failed to handle empreendimento for property %s: %v\n", ext.Codigo, err)
		} else {
			empreendimentoID = empID
		}
	}

	var precoVendaID uint
	if ext.PrecoVenda != nil && ext.PrecoVenda.Ativo {
		pvID, err := is.upsertPrecoVenda(ctx, ext.PrecoVenda)
		if err != nil {
			fmt.Printf("Warning: Failed to handle preco venda for property %s: %v\n", ext.Codigo, err)
		} else {
			precoVendaID = pvID
		}
	}

	var precoAluguelID uint
	if ext.PrecoAluguel != nil && ext.PrecoAluguel.Ativo {
		paID, err := is.upsertPrecoAluguel(ctx, ext.PrecoAluguel)
		if err != nil {
			fmt.Printf("Warning: Failed to handle preco aluguel for property %s: %v\n", ext.Codigo, err)
		} else {
			precoAluguelID = paID
		}
	}

	var corretorPrincipalID uint
	if ext.CorretorPrincipal.Email != "" {
		cpID, err := is.upsertCorretorPrincipal(ctx, &ext.CorretorPrincipal)
		if err != nil {
			fmt.Printf("Warning: Failed to handle corretor principal for property %s: %v\n", ext.Codigo, err)
		} else {
			corretorPrincipalID = cpID
		}
	}

	if isUpdate {
		// Update existing property with new field values AND relationships
		updateReq := &UpdateImovelRequest{
			Titulo:       ext.Titulo,
			Tipo:         ext.Tipo,
			Objetivo:     ext.Objetivo,
			Finalidade:   ext.Finalidade,
			Descricao:    ext.Descricao,
			Metragem:     &ext.Metragem,
			NumQuartos:   &ext.NumQuartos,
			NumSuites:    &ext.NumSuites,
			NumBanheiros: &ext.NumBanheiros,
			NumVagas:     &ext.NumVagas,
			NumAndar:     &ext.NumAndar,
			Unidade:      ext.Unidade,
			Condominio:   &ext.Condominio,
		}

		// Update relationships (use pointers for optional fields)
		if empreendimentoID != 0 {
			updateReq.EmpreendimentoID = &empreendimentoID
		}
		if precoVendaID != 0 {
			updateReq.PrecoVendaID = &precoVendaID
		}
		if precoAluguelID != 0 {
			updateReq.PrecoAluguelID = &precoAluguelID
		}
		if corretorPrincipalID != 0 {
			updateReq.CorretorPrincipalID = &corretorPrincipalID
		}

		imovelResp, err = is.service.UpdateImovel(ctx, imovelID, updateReq)
		if err != nil {
			return nil, fmt.Errorf("failed to update property: %w", err)
		}

		// Update endereco if present
		if ext.Endereco.Rua != "" {
			if err := is.upsertEndereco(ctx, imovelID, &ext.Endereco); err != nil {
				fmt.Printf("Warning: Failed to update endereco for property %s: %v\n", ext.Codigo, err)
			}
		}
	} else {
		// Create endereco first if present
		var enderecoID uint
		if ext.Endereco.Rua != "" {
			enderecoID, err = is.createEndereco(ctx, &ext.Endereco)
			if err != nil {
				return nil, fmt.Errorf("failed to create endereco: %w", err)
			}
		}

		// Create new property with all relationships already upserted above
		createReq := is.transformExternalToCreateRequest(ext, enderecoID, empreendimentoID, precoVendaID, precoAluguelID, corretorPrincipalID)
		imovelResp, err = is.service.CreateImovel(ctx, createReq)
		if err != nil {
			return nil, fmt.Errorf("failed to create property: %w", err)
		}
		imovelID = imovelResp.ID
	}

	// Handle Anexos (Images/Attachments)
	// DELETE old anexos and recreate with current data from external API
	// This ensures removed images are deleted and new images are added
	if err := is.syncAnexosFromImages(ctx, imovelID, ext.Imagens); err != nil {
		fmt.Printf("Warning: Failed to sync attachments for property %s: %v\n", ext.Codigo, err)
	}

	return imovelResp, nil
}

// createEndereco creates a new address and returns its ID
func (is *importService) createEndereco(ctx context.Context, extEndereco *ExternalEndereco) (uint, error) {
	if extEndereco == nil || extEndereco.Rua == "" {
		return 0, fmt.Errorf("endereco is empty")
	}

	endereco := &Endereco{
		Rua:       extEndereco.Rua,
		Numero:    extEndereco.Numero,
		Bairro:    extEndereco.Bairro,
		Cidade:    extEndereco.Cidade,
		Estado:    extEndereco.Estado,
		CEP:       extEndereco.CEP,
		Latitude:  extEndereco.Latitude,
		Longitude: extEndereco.Longitude,
	}

	if err := is.service.CreateEndereco(ctx, endereco); err != nil {
		return 0, fmt.Errorf("failed to create endereco: %w", err)
	}

	return endereco.ID, nil
}

// upsertEndereco creates or updates an address and attaches it to the imovel
func (is *importService) upsertEndereco(ctx context.Context, imovelID uint, extEndereco *ExternalEndereco) error {
	enderecoID, err := is.createEndereco(ctx, extEndereco)
	if err != nil {
		return err
	}

	// Attach the new endereco to the imovel
	return is.service.AttachEndereco(ctx, imovelID, enderecoID)
}

// upsertEmpreendimento creates or updates an enterprise and its nested relationships
func (is *importService) upsertEmpreendimento(ctx context.Context, ext *ExternalEmpreendimento) (uint, error) {
	if ext == nil {
		return 0, fmt.Errorf("empreendimento is nil")
	}

	if ext.ID == 0 {
		return 0, fmt.Errorf("empreendimento has no valid external ID")
	}

	idIntegracao := fmt.Sprintf("%d", ext.ID)

	// Check if empreendimento with this external ID already exists
	var existing Empreendimento
	err := is.service.(*service).repo.(*repository).db.
		Where("id_integracao = ?", idIntegracao).
		First(&existing).Error

	if err == nil {
		// Empreendimento exists, update relevant fields only (skip dates, createdAt)
		updates := map[string]interface{}{
			"titulo":      ext.Titulo,
			"descricao":   ext.Descricao,
			"tipo":        ext.Tipo,
			"status":      ext.Status,
			"localizacao": ext.Localizacao,
		}

		if ext.Finalidade != "" {
			updates["finalidade"] = ext.Finalidade
		}

		// Only update if there are changes (GORM will handle this efficiently)
		if err := is.service.(*service).repo.(*repository).db.
			Model(&existing).
			Updates(updates).Error; err != nil {
			return 0, fmt.Errorf("failed to update empreendimento: %w", err)
		}

		return existing.ID, nil
	}

	// Create new empreendimento - skip fields with date type that cause empty string errors
	empreendimento := &Empreendimento{
		IdIntegracao: idIntegracao,
		Titulo:       ext.Titulo,
		Descricao:    ext.Descricao,
		Tipo:         ext.Tipo,
		Status:       ext.Status,
		Localizacao:  ext.Localizacao,
	}

	if ext.Finalidade != "" {
		empreendimento.Finalidade = ext.Finalidade
	}

	// Use Select to omit problematic fields (data_entrega, etapa_lancamento, endereco_id)
	if err := is.service.(*service).repo.(*repository).db.
		Omit("DataEntrega", "EtapaLancamento", "EnderecoID").
		Create(empreendimento).Error; err != nil {
		return 0, fmt.Errorf("failed to create empreendimento: %w", err)
	}

	return empreendimento.ID, nil
}

// upsertPrecoVenda creates or updates a selling price record
func (is *importService) upsertPrecoVenda(ctx context.Context, ext *ExternalPrecoVenda) (uint, error) {
	if ext == nil {
		return 0, fmt.Errorf("preco venda is nil")
	}

	if ext.ID == 0 {
		return 0, fmt.Errorf("preco venda has no valid external ID")
	}

	idIntegracao := fmt.Sprintf("%d", ext.ID)

	// Check if preco venda with this external ID already exists
	var existing PrecoVenda
	err := is.service.(*service).repo.(*repository).db.
		Where("id_integracao = ?", idIntegracao).
		First(&existing).Error

	if err == nil {
		// Preco venda exists, update it and return its local ID
		existing.Preco = ext.Preco
		existing.AceitaFinanciamentoBancario = ext.AceitaFinanciamentoBancario
		existing.AceitaFinanciamentoDireto = ext.AceitaFinanciamentoDireto
		existing.AceitaPermuta = ext.AceitaPermuta
		existing.AceitaCartaDeCredito = ext.AceitaCartaDeCredito
		existing.AceitaFGTS = ext.AceitaFGTS
		existing.Ativo = ext.Ativo

		if err := is.service.(*service).repo.(*repository).db.Save(&existing).Error; err != nil {
			return 0, fmt.Errorf("failed to update preco venda: %w", err)
		}

		return existing.ID, nil
	}

	// Create new preco venda
	precoVenda := &PrecoVenda{
		IdIntegracao:                idIntegracao,
		Preco:                       ext.Preco,
		AceitaFinanciamentoBancario: ext.AceitaFinanciamentoBancario,
		AceitaFinanciamentoDireto:   ext.AceitaFinanciamentoDireto,
		AceitaPermuta:               ext.AceitaPermuta,
		AceitaCartaDeCredito:        ext.AceitaCartaDeCredito,
		AceitaFGTS:                  ext.AceitaFGTS,
		Ativo:                       ext.Ativo,
	}

	if err := is.service.(*service).repo.(*repository).db.Create(precoVenda).Error; err != nil {
		return 0, fmt.Errorf("failed to create preco venda: %w", err)
	}

	return precoVenda.ID, nil
}

// upsertPrecoAluguel creates or updates a rental price record
func (is *importService) upsertPrecoAluguel(ctx context.Context, ext *ExternalPrecoAluguel) (uint, error) {
	if ext == nil {
		return 0, fmt.Errorf("preco aluguel is nil")
	}

	if ext.ID == 0 {
		return 0, fmt.Errorf("preco aluguel has no valid external ID")
	}

	idIntegracao := fmt.Sprintf("%d", ext.ID)

	// Check if preco aluguel with this external ID already exists
	var existing PrecoAluguel
	err := is.service.(*service).repo.(*repository).db.
		Where("id_integracao = ?", idIntegracao).
		First(&existing).Error

	if err == nil {
		// Preco aluguel exists, update it and return its local ID
		existing.Preco = ext.Preco
		existing.AceitaFiador = ext.AceitaFiador
		existing.Ativo = ext.Ativo

		if err := is.service.(*service).repo.(*repository).db.Save(&existing).Error; err != nil {
			return 0, fmt.Errorf("failed to update preco aluguel: %w", err)
		}

		return existing.ID, nil
	}

	// Create new preco aluguel
	precoAluguel := &PrecoAluguel{
		IdIntegracao: idIntegracao,
		Preco:        ext.Preco,
		AceitaFiador: ext.AceitaFiador,
		Ativo:        ext.Ativo,
	}

	if err := is.service.(*service).repo.(*repository).db.Create(precoAluguel).Error; err != nil {
		return 0, fmt.Errorf("failed to create preco aluguel: %w", err)
	}

	return precoAluguel.ID, nil
}

// setHeaders adds required API headers to the request
func (is *importService) setHeaders(req *http.Request) {
	req.Header.Set("x-api-key", is.apiKey)
	req.Header.Set("x-integration-source", is.integrationSource)
	req.Header.Set("Content-Type", "application/json")
}

// fetchPublishedList fetches the list of published properties
func (is *importService) fetchPublishedList(ctx context.Context, url string) ([]ExternalImovel, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	is.setHeaders(req)

	resp, err := is.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch properties: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp ExternalAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return apiResp.Results.Entities, nil
}

// transformExternalToCreateRequest converts external API response to CreateImovelRequest
func (is *importService) transformExternalToCreateRequest(ext *ExternalDetailedImovel, enderecoID uint, empreendimentoID uint, precoVendaID uint, precoAluguelID uint, corretorPrincipalID uint) *CreateImovelRequest {
	// Default values
	descricao := ext.Descricao
	if descricao == "" {
		descricao = fmt.Sprintf("%s - %s", ext.Titulo, ext.Tipo)
	}

	req := &CreateImovelRequest{
		IdIntegracao:        fmt.Sprintf("%d", ext.ID),
		Titulo:              ext.Titulo,
		Codigo:              ext.Codigo,
		Tipo:                ext.Tipo,
		Objetivo:            ext.Objetivo,
		Finalidade:          ext.Finalidade,
		Descricao:           descricao,
		Metragem:            ext.Metragem,
		NumQuartos:          ext.NumQuartos,
		NumSuites:           ext.NumSuites,
		NumBanheiros:        ext.NumBanheiros,
		NumVagas:            ext.NumVagas,
		NumAndar:            ext.NumAndar,
		Unidade:             ext.Unidade,
		Condominio:          ext.Condominio,
		EnderecoID:          enderecoID,
		CorretorPrincipalID: corretorPrincipalID,
	}

	// Only set empreendimento if it was created successfully
	if empreendimentoID != 0 {
		req.EmpreendimentoID = empreendimentoID
	}

	// Set pricing IDs if they were created successfully
	if precoVendaID != 0 {
		req.PrecoVendaID = precoVendaID
	}

	if precoAluguelID != 0 {
		req.PrecoAluguelID = precoAluguelID
	}

	return req
}

// upsertOrganizacao creates or updates organizacao and returns its ID
func (is *importService) upsertOrganizacao(ctx context.Context, extOrg *ExternalOrganizacao) (uint, error) {
	if extOrg == nil || extOrg.Nome == "" {
		return 0, fmt.Errorf("organizacao is empty")
	}

	// Try to find existing organizacao by external ID
	var org Organizacao

	// Since we don't have IdIntegracao in Organizacao model, we search by Nome
	// This assumes Nome is unique for organizations
	result := is.service.(*service).repo.(*repository).db.Where("nome = ?", extOrg.Nome).First(&org)

	if result.Error == nil {
		// Organizacao exists, update if needed
		if org.Perfil != extOrg.Perfil {
			org.Perfil = extOrg.Perfil
			if err := is.service.(*service).repo.(*repository).db.Save(&org).Error; err != nil {
				return 0, fmt.Errorf("failed to update organizacao: %w", err)
			}
		}
		return org.ID, nil
	}

	// Create new organizacao
	org = Organizacao{
		Nome:   extOrg.Nome,
		Perfil: extOrg.Perfil,
	}

	if err := is.service.(*service).repo.(*repository).db.Create(&org).Error; err != nil {
		return 0, fmt.Errorf("failed to create organizacao: %w", err)
	}

	return org.ID, nil
}

// upsertCorretorPrincipal creates or updates corretor principal and returns its ID
func (is *importService) upsertCorretorPrincipal(ctx context.Context, extCorretor *ExternalCorretor) (uint, error) {
	if extCorretor == nil || extCorretor.Email == "" {
		return 0, fmt.Errorf("corretor principal is empty")
	}

	// First, upsert organizacao
	var organizacaoID uint
	if extCorretor.Organizacao.Nome != "" {
		orgID, err := is.upsertOrganizacao(ctx, &extCorretor.Organizacao)
		if err != nil {
			return 0, fmt.Errorf("failed to upsert organizacao: %w", err)
		}
		organizacaoID = orgID
	}

	// Try to find existing corretor by IdIntegracao
	var corretor CorretorPrincipal
	idIntegracao := fmt.Sprintf("%d", extCorretor.ID)

	result := is.service.(*service).repo.(*repository).db.Where("id_integracao = ?", idIntegracao).First(&corretor)

	if result.Error == nil {
		// Corretor exists, update if needed
		updated := false
		if corretor.Nome != extCorretor.Nome {
			corretor.Nome = extCorretor.Nome
			updated = true
		}
		if corretor.Email != extCorretor.Email {
			corretor.Email = extCorretor.Email
			updated = true
		}
		if corretor.Whatsapp != extCorretor.Whatsapp {
			corretor.Whatsapp = extCorretor.Whatsapp
			updated = true
		}
		if organizacaoID != 0 && corretor.OrganizacaoID != organizacaoID {
			corretor.OrganizacaoID = organizacaoID
			updated = true
		}

		if updated {
			if err := is.service.(*service).repo.(*repository).db.Save(&corretor).Error; err != nil {
				return 0, fmt.Errorf("failed to update corretor principal: %w", err)
			}
		}
		return corretor.ID, nil
	}

	// Create new corretor principal
	corretor = CorretorPrincipal{
		IdIntegracao:   idIntegracao,
		Nome:           extCorretor.Nome,
		Email:          extCorretor.Email,
		Whatsapp:       extCorretor.Whatsapp,
		Idiomas:        extCorretor.Idiomas,
		BairrosAtuacao: extCorretor.BairrosAtuacao,
		OrganizacaoID:  organizacaoID,
	}

	// Don't set FotoID - it will be NULL by default (uint zero value causes FK violation)
	if err := is.service.(*service).repo.(*repository).db.Omit("FotoID").Create(&corretor).Error; err != nil {
		return 0, fmt.Errorf("failed to create corretor principal: %w", err)
	}

	return corretor.ID, nil
}

// addAnexosFromImages adds image attachments to a property
func (is *importService) addAnexosFromImages(ctx context.Context, imovelID uint, imageURLs []string) error {
	// Get existing anexos for this property
	existingAnexos, err := is.service.GetAnexos(ctx, imovelID)
	if err != nil {
		// If error getting existing anexos, log but continue with creation
		fmt.Printf("Warning: Failed to get existing anexos: %v\n", err)
	}

	// Build map of existing URLs for quick lookup
	existingURLs := make(map[string]bool)
	for _, anexo := range existingAnexos {
		existingURLs[anexo.URL] = true
	}

	// Only create anexos that don't already exist
	for i, imageURL := range imageURLs {
		// Skip if this URL already exists
		if existingURLs[imageURL] {
			continue
		}

		anexo := &Anexo{
			Nome:          fmt.Sprintf("Image %d", i+1),
			URL:           imageURL,
			Tipo:          "image",
			Image:         true,
			Video:         false,
			IsExternalURL: true,
			CanPublish:    true,
		}

		if err := is.service.AddAnexo(ctx, imovelID, anexo); err != nil {
			return fmt.Errorf("failed to add image %d: %w", i+1, err)
		}
	}

	return nil
}

// syncAnexosFromImages synchronizes image attachments for a property
// Deletes all existing anexos for this property and recreates them from current external API data
// This ensures that removed images are deleted and new images are added correctly
func (is *importService) syncAnexosFromImages(ctx context.Context, imovelID uint, imageURLs []string) error {
	// Step 1: Delete all existing anexos for this property
	// This ensures removed images from external API are also removed locally
	db := is.service.(*service).repo.(*repository).db
	if err := db.Where("imovel_id = ?", imovelID).Delete(&Anexo{}).Error; err != nil {
		return fmt.Errorf("failed to delete existing anexos: %w", err)
	}

	// Step 2: Create new anexos from current external API data
	for i, imageURL := range imageURLs {
		anexo := &Anexo{
			Nome:          fmt.Sprintf("Image %d", i+1),
			URL:           imageURL,
			Tipo:          "image",
			Image:         true,
			Video:         false,
			IsExternalURL: true,
			CanPublish:    true,
		}

		if err := is.service.AddAnexo(ctx, imovelID, anexo); err != nil {
			return fmt.Errorf("failed to add image %d: %w", i+1, err)
		}
	}

	fmt.Printf("Synced %d anexos for property ID %d\n", len(imageURLs), imovelID)
	return nil
}
