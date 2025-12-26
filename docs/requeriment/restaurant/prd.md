# Restaurant Domain – Requirements Spec (v1)

<role>
Você é um engenheiro de software sênior especialista em Go (Golang).
Vai implementar um backend minimalista porém sólido para o domínio de Restaurantes, priorizando regras de negócio fortes e Aggregate Roots.
</role>

<domain_models>
// Aggregate Root
type Restaurant struct {
    ID                 uuid.UUID
    Name               string
    Slug               string  // "pizza-do-joao" (Unique)
    Description        string
    Status             string  // "DRAFT", "OPEN", "CLOSED", "SUSPENDED"
    Category           string  // "Pizza", "Burgers"
    Rating             int     // 0, 1, 2, 3, 4 ou 5
    TotalReviews       int     // Default 0
    IsOpen             bool    // Campo computado ou flag de override
    DeliveryFee        int64   // unidades monetárias
    MinOrderValue      int64   // unidades monetárias
    PreparationTimeMin int     // em minutos
    SupportsPickup     bool
    SupportsDelivery   bool
    LogoURL            string  // Não obrigatório 
    BannerURL          string  // Não obrigatório 
    CreatedAt      time.Time
    UpdatedAt      time.Time
    
    // Relacionamentos (Carregados com o Aggregate)
    Address        *Address
    OpeningHours   []OpeningHour
    PaymentMethods []PaymentMethod
}

type Address struct {
    ID           uuid.UUID
    RestaurantID uuid.UUID
    Street       string
    Number       string
    Complement   string
    City         string
    State        string // char(2)
    ZipCode      string
    Lat          float64
    Lng          float64
}

type OpeningHour struct {
    ID           uuid.UUID
    RestaurantID uuid.UUID
    Weekday      int // 0=Domingo, 1=Segunda ... 6=Sábado
    OpensAt      int // Minutos a partir da meia-noite (ex: 480 = 08:00)
    ClosesAt     int // Minutos a partir da meia-noite (ex: 120 = 02:00 do dia seguinte)
}

type PaymentMethod struct {
    ID           uuid.UUID
    RestaurantID uuid.UUID
    Method       string // "PIX", "CREDIT_CARD", "DEBIT_CARD"
}
</domain_models>

<instructions>
- **Aggregate Root:** O `Restaurant` é a raiz. `Address`, `OpeningHours` e `PaymentMethods` são entidades filhas que devem ser salvas/gerenciadas preferencialmente dentro da transação do restaurante ou em endpoints filhos estritos.
- **Money Pattern:** Valores monetários (`delivery_fee`, `min_order_value`) devem ser persistidos como `int64` (unidade monetárias) no Go, mapeados para `INTEGER` no banco.
- **Slug Generation:** O Slug deve ser gerado automaticamente a partir do nome se não enviado, garantindo unicidade.
- **Validação:** Use validação manual no UseCase.
</instructions>

<endpoints>
Método  Rota                            Caso de Uso
POST    /restaurants                    CreateRestaurant (Básico + Endereço)
GET     /restaurants                    ListRestaurants
GET     /restaurants/{slug}             GetRestaurantBySlug
PATCH   /restaurants/{id}/open          OpenRestaurant
PATCH   /restaurants/{id}/close         CloseRestaurant
PUT     /restaurants/{id}/hours         UpdateOpeningHours
PUT     /restaurants/{id}/payments      UpdatePaymentMethods
</endpoints>

<business_rules>
R1 - Identidade & Status
- Restaurante nasce sempre com status `DRAFT`.
- `Slug` deve ser único no sistema.

R2 - Operacional
- `IsOpen` é uma combinação de: Status == OPEN + Horário Atual dentro de `restaurant_opening_hours`.
- O endpoint `OpenRestaurant` altera a intenção do lojista (Status=OPEN), mas o frontend deve verificar também os horários.
- Taxa de entrega e Pedido mínimo não podem ser negativos.

R3 - Geolocalização
- Latitude e Longitude são obrigatórios para ativação (Status OPEN), mas opcionais no DRAFT.

R4 - Avaliação
- `Rating` e `TotalReviews` são Read-Only via API de restaurantes (atualizados via triggers ou serviço de Reviews).
</business_rules>

<acceptance_criteria>
**1. Ciclo de Vida & Onboarding**
- [ ] **Criação Default:** Ao criar um restaurante apenas com Nome, o status deve ser `DRAFT`, o `slug` deve ser gerado automaticamente (ex: "Burgers King" -> "burgers-king") e `is_open` deve ser `false`.
- [ ] **Unicidade de Slug:** Tentar criar um segundo restaurante com nome idêntico deve falhar com erro 409 (Conflict).

**2. Gestão de Horários (Integer Minutes)**
- [ ] **Conversão/Validação:** O sistema deve rejeitar valores de minutos < 0 ou > 1439 (23:59).
- [ ] **Madrugada:** O sistema DEVE aceitar `closes_at < opens_at` (ex: Abre 22:00/1320min e Fecha 02:00/120min), interpretando isso como fechamento no dia seguinte.
- [ ] **Colisão:** O sistema deve impedir cadastrar dois intervalos que se sobreponham no mesmo dia da semana.

**3. Regras de Transição (Open/Close)**
- [ ] **Trava de Endereço:** Tentar mudar status para `OPEN` sem ter um `Address` vinculado deve falhar (Erro 400/422).
- [ ] **Trava de Pagamento:** Tentar mudar status para `OPEN` sem nenhum método de pagamento configurado deve falhar.
- [ ] **Trava de Horário:** Tentar mudar status para `OPEN` sem nenhum horário de funcionamento cadastrado deve falhar.
- [ ] **Campo Computado `is_open`:** Ao listar restaurantes (GET), o campo `is_open` no JSON deve ser `true` APENAS SE: (Status == OPEN) **E** (Horário atual estiver dentro de um intervalo válido).

**4. Geolocalização**
- [ ] **Lat/Long:** O endereço deve persistir Latitude e Longitude com precisão decimal correta.
- [ ] **Update de Endereço:** Atualizar o endereço de um restaurante `OPEN` não deve quebrar a consistência, mas deve atualizar os dados imediatamente.
</acceptance_criteria>

<errors>
Código  Quando
400     Dados inválidos (Slug malformado, CEP inválido).
404     Restaurante não encontrado.
409     Conflito (Slug já existe).
500     Erro interno.
</errors>

