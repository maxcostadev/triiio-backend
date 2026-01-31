package imoveis

// External API DTOs for mapping responses from dev-api-backend.pi8.com.br

// ExternalAPIResponse represents the top-level response structure
type ExternalAPIResponse struct {
	Results ExternalResults `json:"results"`
}

// ExternalResults contains the entities array
type ExternalResults struct {
	Entities []ExternalImovel `json:"entities"`
}

// ExternalImovel represents a property from the external API
type ExternalImovel struct {
	ID                uint                  `json:"id"`
	Codigo            string                `json:"codigo"`
	Titulo            string                `json:"titulo"`
	Tipo              string                `json:"tipo"`       // APARTAMENTO, CASA, COMERCIAL, etc
	Objetivo          string                `json:"objetivo"`   // VENDER, ALUGAR
	Finalidade        string                `json:"finalidade"` // RESIDENCIAL, COMERCIAL
	Metragem          float64               `json:"metragem"`
	NumQuartos        int                   `json:"numQuartos"`
	NumSuites         int                   `json:"numSuites"`
	NumBanheiros      int                   `json:"numBanheiros"`
	NumVagas          int                   `json:"numVagas"`
	NumAndar          int                   `json:"numAndar"`
	Unidade           string                `json:"unidade"`
	Condominio        float64               `json:"condominio"`
	Preco             float64               `json:"preco"`
	Status            string                `json:"status"` // PUBLICADO, EM_EDICAO, ARQUIVADO
	Visualizacoes     int                   `json:"visualizacoes"`
	InfoAnuncio       string                `json:"infoAnuncio"`
	Imagens           []string              `json:"imagens"`
	Endereco          ExternalEndereco      `json:"endereco"`
	CorretorPrincipal ExternalCorretor      `json:"corretorPrincipal"`
	PrecoVenda        *ExternalPrecoVenda   `json:"precoVenda"`
	PrecoAluguel      *ExternalPrecoAluguel `json:"precoAluguel"`
	Compartilhamentos []interface{}         `json:"compartilhamentos"`
}

// ExternalEndereco represents address from external API
type ExternalEndereco struct {
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

// ExternalOrganizacao represents organization from external API
type ExternalOrganizacao struct {
	ID     uint   `json:"id"`
	Nome   string `json:"nome"`
	Perfil string `json:"perfil"`
}

// ExternalFoto represents photo from external API
type ExternalFoto struct {
	URL     string `json:"url"`
	Tipo    string `json:"tipo"`
	Tamanho int64  `json:"tamanho"`
}

// ExternalCorretor represents broker from external API
type ExternalCorretor struct {
	ID             uint                `json:"id"`
	Nome           string              `json:"nome"`
	Email          string              `json:"email"`
	Whatsapp       string              `json:"whatsapp"`
	Foto           *ExternalFoto       `json:"foto"`
	Idiomas        []string            `json:"idiomas"`
	BairrosAtuacao []string            `json:"bairrosAtuacao"`
	Organizacao    ExternalOrganizacao `json:"organizacao"`
}

// ExternalPrecoVenda represents selling price from external API
type ExternalPrecoVenda struct {
	ID                          uint           `json:"id"`
	Preco                       float64        `json:"preco"`
	AceitaFinanciamentoBancario bool           `json:"aceitaFinanciamentoBancario"`
	AceitaFinanciamentoDireto   bool           `json:"aceitaFinanciamentoDireto"`
	AceitaPermuta               bool           `json:"aceitaPermuta"`
	AceitaCartaDeCredito        bool           `json:"aceitaCartaDeCredito"`
	AceitaFGTS                  bool           `json:"aceitaFGTS"`
	Ativo                       bool           `json:"ativo"`
	Pacote                      ExternalPacote `json:"pacote"`
}

// ExternalPrecoAluguel represents rental price from external API
type ExternalPrecoAluguel struct {
	ID           uint    `json:"id"`
	Preco        float64 `json:"preco"`
	AceitaFiador bool    `json:"aceitaFiador"`
	Ativo        bool    `json:"ativo"`
}

// ExternalPacote represents package from external API
type ExternalPacote struct {
	ID         uint   `json:"id"`
	Titulo     string `json:"titulo"`
	Descricao  string `json:"descricao"`
	Exclusivo  bool   `json:"exclusivo"`
	EmDestaque bool   `json:"emDestaque"`
}

// ExternalDetailedImovel represents detailed property info from single property endpoint
type ExternalDetailedImovel struct {
	ID                uint                    `json:"id"`
	Codigo            string                  `json:"codigo"`
	Titulo            string                  `json:"titulo"`
	Descricao         string                  `json:"descricao"`
	Tipo              string                  `json:"tipo"`
	Objetivo          string                  `json:"objetivo"`
	Finalidade        string                  `json:"finalidade"`
	Metragem          float64                 `json:"metragem"`
	NumQuartos        int                     `json:"numQuartos"`
	NumSuites         int                     `json:"numSuites"`
	NumBanheiros      int                     `json:"numBanheiros"`
	NumVagas          int                     `json:"numVagas"`
	NumAndar          int                     `json:"numAndar"`
	Unidade           string                  `json:"unidade"`
	Condominio        float64                 `json:"condominio"`
	Status            string                  `json:"status"`
	Visualizacoes     int                     `json:"visualizacoes"`
	Imagens           []string                `json:"imagens"`
	Endereco          ExternalEndereco        `json:"endereco"`
	CorretorPrincipal ExternalCorretor        `json:"corretorPrincipal"`
	PrecoVenda        *ExternalPrecoVenda     `json:"precoVenda"`
	PrecoAluguel      *ExternalPrecoAluguel   `json:"precoAluguel"`
	Empreendimento    *ExternalEmpreendimento `json:"empreendimento"`
}

// ExternalEmpreendimento represents enterprise from external API
type ExternalEmpreendimento struct {
	ID              uint             `json:"id"`
	Codigo          string           `json:"codigo"`
	Titulo          string           `json:"titulo"`
	Descricao       string           `json:"descricao"`
	DataEntrega     string           `json:"data_entrega"`
	EtapaLancamento string           `json:"etapa_lancamento"`
	Finalidade      string           `json:"finalidade"`
	Tipo            string           `json:"tipo"`
	Status          string           `json:"status"`
	Localizacao     string           `json:"localizacao"`
	Endereco        ExternalEndereco `json:"endereco"`
	Torres          []ExternalTorre  `json:"torres"`
	Plantas         []ExternalPlanta `json:"plantas"`
}

// ExternalTorre represents tower from external API
type ExternalTorre struct {
	ID              uint   `json:"id"`
	Nome            string `json:"nome"`
	TotalColunas    int    `json:"totalColunas"`
	TotalElevadores int    `json:"totalElevadores"`
	TotalPavimentos int    `json:"totalPavimentos"`
	TotalUnidades   int    `json:"totalUnidades"`
}

// ExternalPlanta represents floor plan from external API
type ExternalPlanta struct {
	ID       uint     `json:"id"`
	Nome     string   `json:"nome"`
	Metragem float64  `json:"metragem"`
	Imagens  []string `json:"imagens"`
}
