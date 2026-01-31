package imoveis

import "time"

// CreateImovelRequest represents property creation request
type CreateImovelRequest struct {
	IdIntegracao  string  `json:"id_integracao" binding:"omitempty"`
	Titulo        string  `json:"titulo" binding:"required,min=3,max=255"`
	Codigo        string  `json:"codigo" binding:"required,min=1,max=50"`
	Tipo          string  `json:"tipo" binding:"required,oneof=APARTAMENTO CASA COMERCIAL SALA_COMERCIAL TERRENO GALPAO"`
	Objetivo      string  `json:"objetivo" binding:"required,oneof=VENDER ALUGAR"`
	Finalidade    string  `json:"finalidade" binding:"required,oneof=RESIDENCIAL COMERCIAL MISTO"`
	Descricao     string  `json:"descricao" binding:"required,min=10,max=5000"`
	Metragem      float64 `json:"metragem" binding:"required,gt=0"`
	NumQuartos    int     `json:"numQuartos" binding:"min=0"`
	NumSuites     int     `json:"numSuites" binding:"min=0"`
	NumBanheiros  int     `json:"numBanheiros" binding:"min=0"`
	NumVagas      int     `json:"numVagas" binding:"min=0"`
	NumAndar      int     `json:"numAndar" binding:"omitempty"`
	Unidade       string  `json:"unidade" binding:"omitempty,max=20"`
	Condominio    float64 `json:"condominio" binding:"min=0"`
	IPTU          float64 `json:"iptu" binding:"min=0"`
	InscricaoIPTU string  `json:"inscricaoIPTU" binding:"omitempty,max=50"`

	// Relations
	EnderecoID          uint   `json:"endereco_id" binding:"required"`
	EmpreendimentoID    uint   `json:"empreendimento_id" binding:"omitempty"`
	PlantaID            uint   `json:"planta_id" binding:"omitempty"`
	CorretorPrincipalID uint   `json:"corretor_principal_id" binding:"omitempty"`
	PacoteID            uint   `json:"pacote_id" binding:"omitempty"`
	PrecoVendaID        uint   `json:"preco_venda_id" binding:"omitempty"`
	PrecoAluguelID      uint   `json:"preco_aluguel_id" binding:"omitempty"`
	Caracteristicas     []uint `json:"caracteristicas" binding:"omitempty,dive"`
}

// UpdateImovelRequest represents property update request
type UpdateImovelRequest struct {
	Titulo        string   `json:"titulo" binding:"omitempty,min=3,max=255"`
	Codigo        string   `json:"codigo" binding:"omitempty,min=1,max=50"`
	Tipo          string   `json:"tipo" binding:"omitempty,oneof=APARTAMENTO CASA COMERCIAL SALA_COMERCIAL TERRENO GALPAO"`
	Objetivo      string   `json:"objetivo" binding:"omitempty,oneof=VENDER ALUGAR"`
	Finalidade    string   `json:"finalidade" binding:"omitempty,oneof=RESIDENCIAL COMERCIAL MISTO"`
	Descricao     string   `json:"descricao" binding:"omitempty,min=10,max=5000"`
	Metragem      *float64 `json:"metragem" binding:"omitempty,gt=0"`
	NumQuartos    *int     `json:"numQuartos" binding:"omitempty,min=0"`
	NumSuites     *int     `json:"numSuites" binding:"omitempty,min=0"`
	NumBanheiros  *int     `json:"numBanheiros" binding:"omitempty,min=0"`
	NumVagas      *int     `json:"numVagas" binding:"omitempty,min=0"`
	NumAndar      *int     `json:"numAndar" binding:"omitempty"`
	Unidade       string   `json:"unidade" binding:"omitempty,max=20"`
	Condominio    *float64 `json:"condominio" binding:"omitempty,min=0"`
	IPTU          *float64 `json:"iptu" binding:"omitempty,min=0"`
	InscricaoIPTU string   `json:"inscricaoIPTU" binding:"omitempty,max=50"`

	// Relations
	EnderecoID          *uint  `json:"endereco_id" binding:"omitempty"`
	EmpreendimentoID    *uint  `json:"empreendimento_id" binding:"omitempty"`
	PlantaID            *uint  `json:"planta_id" binding:"omitempty"`
	CorretorPrincipalID *uint  `json:"corretor_principal_id" binding:"omitempty"`
	PacoteID            *uint  `json:"pacote_id" binding:"omitempty"`
	PrecoVendaID        *uint  `json:"preco_venda_id" binding:"omitempty"`
	PrecoAluguelID      *uint  `json:"preco_aluguel_id" binding:"omitempty"`
	Status              string `json:"status" binding:"omitempty,oneof=PUBLICADO EM_EDICAO ARQUIVADO"`
	Published           *bool  `json:"published" binding:"omitempty"`
	Closed              *bool  `json:"closed" binding:"omitempty"`
	Caracteristicas     []uint `json:"caracteristicas" binding:"omitempty,dive"`
}

// ImovelResponse represents property response
type ImovelResponse struct {
	ID            uint    `json:"id"`
	IdIntegracao  string  `json:"id_integracao"`
	Titulo        string  `json:"titulo"`
	Codigo        string  `json:"codigo"`
	SeqCodigo     int     `json:"seqCodigo"`
	Tipo          string  `json:"tipo"`
	Objetivo      string  `json:"objetivo"`
	Finalidade    string  `json:"finalidade"`
	Descricao     string  `json:"descricao"`
	Metragem      float64 `json:"metragem"`
	NumQuartos    int     `json:"numQuartos"`
	NumSuites     int     `json:"numSuites"`
	NumBanheiros  int     `json:"numBanheiros"`
	NumVagas      int     `json:"numVagas"`
	NumAndar      int     `json:"numAndar"`
	Unidade       string  `json:"unidade"`
	Condominio    float64 `json:"condominio"`
	IPTU          float64 `json:"iptu"`
	InscricaoIPTU string  `json:"inscricaoIPTU"`

	// Relations
	Endereco          *EnderecoResponse          `json:"endereco,omitempty"`
	Empreendimento    *EmpreendimentoResponse    `json:"empreendimento,omitempty"`
	Planta            *PlantaResponse            `json:"planta,omitempty"`
	CorretorPrincipal *CorretorPrincipalResponse `json:"corretorPrincipal,omitempty"`
	Pacote            *PacoteResponse            `json:"pacote,omitempty"`
	PrecoVenda        *PrecoVendaResponse        `json:"precoVenda,omitempty"`
	PrecoAluguel      *PrecoAluguelResponse      `json:"precoAluguel,omitempty"`
	Anexos            []AnexoResponse            `json:"anexos,omitempty"`
	Caracteristicas   []CaracteristicaResponse   `json:"caracteristicas,omitempty"`

	// Metadata
	Status        string    `json:"status"`
	Published     bool      `json:"published"`
	Closed        bool      `json:"closed"`
	Visualizacoes int       `json:"visualizacoes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AnexoResponse represents attachment response
type AnexoResponse struct {
	ID            uint      `json:"id"`
	Nome          string    `json:"nome"`
	Path          string    `json:"path"`
	Tamanho       int64     `json:"tamanho"`
	Tipo          string    `json:"tipo"`
	URL           string    `json:"url"`
	CanPublish    bool      `json:"canPublish"`
	Image         bool      `json:"image"`
	Video         bool      `json:"video"`
	IsExternalURL bool      `json:"isExternalUrl"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// EnderecoResponse represents address response
type EnderecoResponse struct {
	ID        uint    `json:"id"`
	Rua       string  `json:"rua"`
	Numero    int     `json:"numero"`
	Bairro    string  `json:"bairro"`
	Cidade    string  `json:"cidade"`
	Estado    string  `json:"estado"`
	CEP       string  `json:"cep"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// PlantaResponse represents floor plan response
type PlantaResponse struct {
	ID        uint            `json:"id"`
	Nome      string          `json:"nome"`
	Metragem  float64         `json:"metragem"`
	Anexos    []AnexoResponse `json:"anexos,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// TorresResponse represents tower response
type TorresResponse struct {
	ID              uint      `json:"id"`
	Nome            string    `json:"nome"`
	TotalColunas    int       `json:"totalColunas"`
	TotalElevadores int       `json:"totalElevadores"`
	TotalPavimentos int       `json:"totalPavimentos"`
	TotalUnidades   int       `json:"totalUnidades"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreateEmpreendimentoRequest represents enterprise creation request
type CreateEmpreendimentoRequest struct {
	Titulo          string `json:"titulo" binding:"required,min=3,max=255"`
	Descricao       string `json:"descricao" binding:"required,min=10,max=5000"`
	DataEntrega     string `json:"data_entrega" binding:"omitempty,datetime=2006-01-02"`
	EtapaLancamento string `json:"etapa_lancamento" binding:"omitempty,oneof=LANCAMENTO PRE_LANCAMENTO PRONTO EM_CONSTRUCAO"`
	Finalidade      string `json:"finalidade" binding:"omitempty,oneof=RESIDENCIAL COMERCIAL MISTO"`
	Tipo            string `json:"tipo" binding:"omitempty,oneof=APARTAMENTO CASA COMERCIAL SALA_COMERCIAL TERRENO GALPAO"`
	Status          string `json:"status" binding:"omitempty,oneof=PUBLICADO EM_EDICAO ARQUIVADO"`
	Localizacao     string `json:"localizacao" binding:"omitempty,max=255"`
	EnderecoID      uint   `json:"endereco_id" binding:"omitempty"`
}

// UpdateEmpreendimentoRequest represents enterprise update request
type UpdateEmpreendimentoRequest struct {
	Titulo          string `json:"titulo" binding:"omitempty,min=3,max=255"`
	Descricao       string `json:"descricao" binding:"omitempty,min=10,max=5000"`
	DataEntrega     string `json:"data_entrega" binding:"omitempty,datetime=2006-01-02"`
	EtapaLancamento string `json:"etapa_lancamento" binding:"omitempty,oneof=LANCAMENTO PRE_LANCAMENTO PRONTO EM_CONSTRUCAO"`
	Finalidade      string `json:"finalidade" binding:"omitempty,oneof=RESIDENCIAL COMERCIAL MISTO"`
	Tipo            string `json:"tipo" binding:"omitempty,oneof=APARTAMENTO CASA COMERCIAL SALA_COMERCIAL TERRENO GALPAO"`
	Status          string `json:"status" binding:"omitempty,oneof=PUBLICADO EM_EDICAO ARQUIVADO"`
	Localizacao     string `json:"localizacao" binding:"omitempty,max=255"`
	EnderecoID      *uint  `json:"endereco_id" binding:"omitempty"`
}

// EmpreendimentoResponse represents enterprise response
type EmpreendimentoResponse struct {
	ID              uint                     `json:"id"`
	Titulo          string                   `json:"titulo"`
	Descricao       string                   `json:"descricao"`
	DataEntrega     string                   `json:"data_entrega,omitempty"`
	EtapaLancamento string                   `json:"etapa_lancamento,omitempty"`
	Finalidade      string                   `json:"finalidade,omitempty"`
	Tipo            string                   `json:"tipo,omitempty"`
	Status          string                   `json:"status,omitempty"`
	Localizacao     string                   `json:"localizacao"`
	Endereco        *EnderecoResponse        `json:"endereco,omitempty"`
	Plantas         []PlantaResponse         `json:"plantas,omitempty"`
	Torres          []TorresResponse         `json:"torres,omitempty"`
	Caracteristicas []CaracteristicaResponse `json:"caracteristicas,omitempty"`
	Anexos          []AnexoResponse          `json:"anexos,omitempty"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

// OrganizacaoResponse represents organization response
type OrganizacaoResponse struct {
	ID     uint   `json:"id"`
	Nome   string `json:"nome"`
	Perfil string `json:"perfil"`
}

// CorretorPrincipalResponse represents real estate agent response
type CorretorPrincipalResponse struct {
	ID             uint                 `json:"id"`
	Nome           string               `json:"nome"`
	Email          string               `json:"email"`
	Whatsapp       string               `json:"whatsapp"`
	Foto           *AnexoResponse       `json:"foto,omitempty"`
	Idiomas        []string             `json:"idiomas"`
	BairrosAtuacao []string             `json:"bairrosAtuacao"`
	Organizacao    *OrganizacaoResponse `json:"organizacao,omitempty"`
}

// PacoteResponse represents package response
type PacoteResponse struct {
	ID         uint      `json:"id"`
	Titulo     string    `json:"titulo"`
	Descricao  string    `json:"descricao"`
	Exclusivo  bool      `json:"exclusivo"`
	EmDestaque bool      `json:"em_destaque"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CaracteristicaResponse represents characteristic response
type CaracteristicaResponse struct {
	ID            uint      `json:"id"`
	Nome          string    `json:"nome"`
	CategoriaID   uint      `json:"categoria_id,omitempty"`
	CategoriaNome string    `json:"categoria_nome,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PrecoVendaResponse represents selling price response
type PrecoVendaResponse struct {
	ID                          uint      `json:"id"`
	Preco                       float64   `json:"preco"`
	AceitaFinanciamentoBancario bool      `json:"aceitaFinanciamentoBancario"`
	AceitaFinanciamentoDireto   bool      `json:"aceitaFinanciamentoDireto"`
	AceitaPermuta               bool      `json:"aceitaPermuta"`
	AceitaCartaDeCredito        bool      `json:"aceitaCartaDeCredito"`
	AceitaFGTS                  bool      `json:"aceitaFGTS"`
	Ativo                       bool      `json:"ativo"`
	PacoteTitulo                string    `json:"pacote_titulo,omitempty"`
	PacoteDescricao             string    `json:"pacote_descricao,omitempty"`
	PacoteExclusivo             bool      `json:"pacote_exclusivo,omitempty"`
	PacoteEmDestaque            bool      `json:"pacote_em_destaque,omitempty"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}

// PrecoAluguelResponse represents rental price response
type PrecoAluguelResponse struct {
	ID           uint      `json:"id"`
	Preco        float64   `json:"preco,omitempty"`
	AceitaFiador bool      `json:"aceitaFiador"`
	Ativo        bool      `json:"ativo"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ImovelListQuery represents query parameters for listing properties
type ImovelListQuery struct {
	Page             int     `form:"page,default=1" binding:"min=1"`
	Limit            int     `form:"limit,default=10" binding:"min=1,max=100"`
	Codigo           string  `form:"codigo" binding:"omitempty,max=50"`
	Tipo             string  `form:"tipo" binding:"omitempty,oneof=APARTAMENTO CASA COMERCIAL SALA_COMERCIAL TERRENO GALPAO"`
	Objetivo         string  `form:"objetivo" binding:"omitempty,oneof=VENDER ALUGAR"`
	Finalidade       string  `form:"finalidade" binding:"omitempty,oneof=RESIDENCIAL COMERCIAL MISTO"`
	Status           string  `form:"status" binding:"omitempty,oneof=PUBLICADO EM_EDICAO ARQUIVADO"`
	Published        *bool   `form:"published" binding:"omitempty"`
	MinPreco         float64 `form:"min_preco" binding:"omitempty,min=0"`
	MaxPreco         float64 `form:"max_preco" binding:"omitempty,min=0"`
	MinMetragem      float64 `form:"min_metragem" binding:"omitempty,min=0"`
	MaxMetragem      float64 `form:"max_metragem" binding:"omitempty,min=0"`
	Rua              string  `form:"rua" binding:"omitempty,max=200"`
	Cidade           string  `form:"cidade" binding:"omitempty,max=100"`
	Bairro           string  `form:"bairro" binding:"omitempty,max=100"`
	NumQuartos       int     `form:"num_quartos" binding:"omitempty,min=0"`
	NumBanheiros     int     `form:"num_banheiros" binding:"omitempty,min=0"`
	NumGaragens      int     `form:"num_garagens" binding:"omitempty,min=0"`
	EmpreendimentoID uint    `form:"empreendimento_id" binding:"omitempty"`
	Sort             string  `form:"sort" binding:"omitempty,oneof=created_at updated_at preco titulo metragem"`
	Order            string  `form:"order,default=desc" binding:"oneof=asc desc"`
}

// ImovelListResponse represents paginated property list response
type ImovelListResponse struct {
	Total   int64            `json:"total"`
	Page    int              `json:"page"`
	Limit   int              `json:"limit"`
	Pages   int64            `json:"pages"`
	HasNext bool             `json:"hasNext"`
	HasPrev bool             `json:"hasPrev"`
	Results []ImovelResponse `json:"results"`
}
