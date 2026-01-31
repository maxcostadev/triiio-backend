package imoveis

import (
	"time"

	"gorm.io/gorm"
)

// Anexo represents an attachment (image, video, etc.)
type Anexo struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	Nome             string         `json:"nome"`
	Path             string         `json:"path"`
	Tamanho          int64          `json:"tamanho"`
	Tipo             string         `json:"tipo"`
	URL              string         `json:"url"`
	CanPublish       bool           `json:"canPublish"`
	Image            bool           `json:"image"`
	Video            bool           `json:"video"`
	IsExternalURL    bool           `json:"isExternalUrl"`
	ImovelID         *uint          `json:"imovel_id,omitempty"`
	EmpreendimentoID *uint          `json:"empreendimento_id,omitempty"`
	PlantaID         *uint          `json:"planta_id,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// Endereco represents an address
type Endereco struct {
	ID        uint    `gorm:"primarykey" json:"id"`
	Rua       string  `json:"rua"`
	Numero    int     `json:"numero"`
	Bairro    string  `json:"bairro"`
	Cidade    string  `json:"cidade"`
	Estado    string  `json:"estado"`
	CEP       string  `json:"cep"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Plantas struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	Nome             string         `json:"nome"`
	Metragem         float64        `json:"metragem"`
	EmpreendimentoID uint           `json:"empreendimento_id,omitempty"`
	Anexos           []Anexo        `gorm:"foreignKey:PlantaID" json:"anexos,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

type Organizacao struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Nome      string         `json:"nome"`
	Perfil    string         `json:"perfil"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the table name used by GORM (prevents using "organizacaos")
func (Organizacao) TableName() string {
	return "organizacoes"
}

// CorretorPrincipal represents a real estate agent
type CorretorPrincipal struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	IdIntegracao   string         `gorm:"uniqueIndex" json:"id_integracao,omitempty"`
	Nome           string         `json:"nome"`
	Email          string         `json:"email"`
	Whatsapp       string         `json:"whatsapp"`
	FotoID         uint           `json:"foto_id,omitempty"`
	Foto           *Anexo         `gorm:"foreignKey:FotoID" json:"foto,omitempty"`
	Idiomas        []string       `gorm:"type:text[]" json:"idiomas"`
	BairrosAtuacao []string       `gorm:"type:text[]" json:"bairros_atuacao"`
	OrganizacaoID  uint           `json:"organizacao_id"`
	Organizacao    *Organizacao   `gorm:"foreignKey:OrganizacaoID" json:"organizacao,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the table name used by GORM (prevents using "corretor_principals")
func (CorretorPrincipal) TableName() string {
	return "corretores_principais"
}

type Pacote struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	IdIntegracao string         `gorm:"uniqueIndex" json:"id_integracao,omitempty"`
	Titulo       string         `json:"titulo"`
	Descricao    string         `json:"descricao"`
	Exclusivo    bool           `json:"exclusivo"`
	EmDestaque   bool           `json:"em_destaque"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Caracteristica represents a feature/characteristic
type Caracteristica struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Nome          string         `json:"nome"`
	CategoriaID   uint           `json:"categoria_id,omitempty"`
	CategoriaNome string         `json:"categoria_nome,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Empreendimento struct {
	ID              uint             `gorm:"primarykey" json:"id"`
	IdIntegracao    string           `gorm:"uniqueIndex" json:"id_integracao,omitempty"`
	Titulo          string           `json:"titulo"`
	Descricao       string           `json:"descricao"`
	DataEntrega     string           `json:"data_entrega,omitempty"`
	EtapaLancamento string           `json:"etapa_lancamento,omitempty"`
	Finalidade      string           `json:"finalidade,omitempty"`
	Tipo            string           `json:"tipo,omitempty"`
	Status          string           `json:"status,omitempty"`
	Localizacao     string           `json:"localizacao"`
	EnderecoID      uint             `json:"endereco_id,omitempty"`
	Endereco        *Endereco        `gorm:"foreignKey:EnderecoID" json:"endereco,omitempty"`
	Plantas         []Plantas        `gorm:"foreignKey:EmpreendimentoID" json:"plantas,omitempty"`
	Torres          []Torres         `gorm:"foreignKey:EmpreendimentoID" json:"torres,omitempty"`
	Caracteristicas []Caracteristica `gorm:"many2many:empreendimento_caracteristicas;" json:"caracteristicas,omitempty"`
	Anexos          []Anexo          `gorm:"foreignKey:EmpreendimentoID" json:"anexos,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `gorm:"index" json:"-"`
}

type Torres struct {
	ID               uint            `gorm:"primarykey" json:"id"`
	Nome             string          `json:"nome"`
	TotalColunas     int             `json:"totalColunas"`
	TotalElevadores  int             `json:"totalElevadores"`
	TotalPavimentos  int             `json:"totalPavimentos"`
	TotalUnidades    int             `json:"totalUnidades"`
	EmpreendimentoID uint            `json:"empreendimento_id"`
	Empreendimento   *Empreendimento `gorm:"foreignKey:EmpreendimentoID" json:"empreendimento,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"-"`
}

// PrecoVenda represents selling price details
type PrecoVenda struct {
	ID                          uint           `gorm:"primarykey" json:"id"`
	IdIntegracao                string         `gorm:"uniqueIndex" json:"id_integracao,omitempty"`
	Preco                       float64        `json:"preco"`
	AceitaFinanciamentoBancario bool           `json:"aceitaFinanciamentoBancario"`
	AceitaFinanciamentoDireto   bool           `json:"aceitaFinanciamentoDireto"`
	AceitaPermuta               bool           `json:"aceitaPermuta"`
	AceitaCartaDeCredito        bool           `json:"aceitaCartaDeCredito"`
	AceitaFGTS                  bool           `json:"aceitaFGTS"`
	Ativo                       bool           `json:"ativo"`
	PacoteTitulo                string         `json:"pacote_titulo,omitempty"`
	PacoteDescricao             string         `json:"pacote_descricao,omitempty"`
	PacoteExclusivo             bool           `json:"pacote_exclusivo,omitempty"`
	PacoteEmDestaque            bool           `json:"pacote_em_destaque,omitempty"`
	CreatedAt                   time.Time      `json:"created_at"`
	UpdatedAt                   time.Time      `json:"updated_at"`
	DeletedAt                   gorm.DeletedAt `gorm:"index" json:"-"`
}

// PrecoAluguel represents rental price details
type PrecoAluguel struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	IdIntegracao string         `gorm:"uniqueIndex" json:"id_integracao,omitempty"`
	Preco        float64        `json:"preco,omitempty"`
	AceitaFiador bool           `json:"aceitaFiador"`
	Ativo        bool           `json:"ativo"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Imovel represents a real estate property
type Imovel struct {
	ID            uint   `gorm:"primarykey" json:"id"`
	Id_Integracao string `gorm:"uniqueIndex;not null" json:"id_integracao"`
	Titulo        string `gorm:"not null" json:"titulo"`
	Codigo        string `gorm:"uniqueIndex;not null" json:"codigo"`
	SeqCodigo     int    `json:"seqCodigo"`
	Tipo          string `json:"tipo"`       // APARTAMENTO, CASA, COMERCIAL, etc
	Objetivo      string `json:"objetivo"`   // VENDER, ALUGAR
	Finalidade    string `json:"finalidade"` // RESIDENCIAL, COMERCIAL
	Descricao     string `gorm:"type:text" json:"descricao"`

	// Property Details
	Metragem     float64 `json:"metragem"`
	NumQuartos   int     `json:"numQuartos"`
	NumSuites    int     `json:"numSuites"`
	NumBanheiros int     `json:"numBanheiros"`
	NumVagas     int     `json:"numVagas"`
	NumAndar     int     `json:"numAndar"`
	Unidade      string  `json:"unidade"`

	// Financial Details
	Condominio    float64 `json:"condominio"`
	IPTU          float64 `gorm:"column:iptu" json:"iptu"`
	InscricaoIPTU string  `gorm:"column:inscricao_iptu" json:"inscricaoIPTU"`

	// Location & Address
	EnderecoID uint      `json:"endereco_id"`
	Endereco   *Endereco `gorm:"foreignKey:EnderecoID" json:"endereco"`

	// Enterprise/Empreendimento relation
	EmpreendimentoID uint            `json:"empreendimento_id,omitempty"`
	Empreendimento   *Empreendimento `gorm:"foreignKey:EmpreendimentoID" json:"empreendimento,omitempty"`
	// Pricing
	PrecoVendaID uint        `json:"preco_venda_id,omitempty"`
	PrecoVenda   *PrecoVenda `gorm:"foreignKey:PrecoVendaID" json:"precoVenda"`

	PrecoAluguelID uint          `json:"preco_aluguel_id,omitempty"`
	PrecoAluguel   *PrecoAluguel `gorm:"foreignKey:PrecoAluguelID" json:"precoAluguel"`

	Anexos []Anexo `gorm:"foreignKey:ImovelID" json:"anexos,omitempty"`

	// Status & Publishing
	Status    string `json:"status"` // PUBLICADO, EM_EDICAO, ARQUIVADO
	Published bool   `gorm:"default:false" json:"published"`
	Closed    bool   `gorm:"default:false" json:"closed"`

	// Plant reference
	PlantaID uint     `json:"plantaID,omitempty"`
	Planta   *Plantas `gorm:"foreignKey:PlantaID" json:"planta,omitempty"`

	// Corretor Principal
	CorretorPrincipalID uint               `json:"corretor_principal_id,omitempty"`
	CorretorPrincipal   *CorretorPrincipal `gorm:"foreignKey:CorretorPrincipalID" json:"corretorPrincipal,omitempty"`

	// Package
	PacoteID uint    `json:"pacote_id,omitempty"`
	Pacote   *Pacote `gorm:"foreignKey:PacoteID" json:"pacote,omitempty"`

	// Characteristics
	Caracteristicas []Caracteristica `gorm:"many2many:imovel_caracteristicas;" json:"caracteristicas,omitempty"`

	// Metadata
	Visualizacoes int            `gorm:"default:0" json:"visualizacoes"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name
func (Imovel) TableName() string {
	return "imoveis"
}
