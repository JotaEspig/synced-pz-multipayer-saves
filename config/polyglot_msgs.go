package config

import "github.com/charmbracelet/log"

var dict = make(map[int]map[string]string, LANG_END-LANG_START-2) // -2 because LANG_EN is not included

// GTM is a shortcut for GetTranslatedMsg
func GTM(msg string) string {
	return GetTranslatedMsg(msg)
}

// GetTranslatedMsg receives a message in English and returns it in the language set by the user
func GetTranslatedMsg(msg string) string {
	// Default language is English
	lang := LANG_EN
	if IsLanguageValid(Launguage) {
		lang = Launguage
	}

	if lang == LANG_EN {
		return msg
	}

	dictInLang, ok := dict[lang]
	if !ok {
		log.Fatal("Language not found")
	}
	translatedMsg, ok := dictInLang[msg]
	if !ok {
		log.Error("Translated message not found")
		return msg
	}
	return translatedMsg
}

func init() {
	// Portuguese
	ptbrDict := make(map[string]string)

	ptbrDict["Press any key to exit..."] = "Pressione qualquer tecla para sair..."
	ptbrDict["Press any key to continue..."] = "Pressione qualquer tecla para continuar..."
	ptbrDict["First time setup"] = "Configuração inicial"
	ptbrDict["Type of servers to list"] = "Tipo de servidores para listar"
	ptbrDict["No argument for config"] = "Nenhum foi argumento para config"
	ptbrDict["Usage: "] = "Uso: "
	ptbrDict["  syncedpz help = shows this message"] = "  syncedpz help = mostra esta mensagem"
	ptbrDict["  syncedpz menu = use menu mode"] = "  syncedpz menu = usa o modo menu"
	ptbrDict["  syncedpz config [setup | list] = sets up or list the syncedpz configuration"] = "  syncedpz config [setup | list] = configura ou lista a configuração do syncedpz"
	ptbrDict["  syncedpz list -type [local | synced] = list servers according to its type (default is local))"] = "  syncedpz list -type [local | synced] = lista servidores de acordo com seu tipo (padrão é local))"
	ptbrDict["  syncedpz add = adds a new synced PZ server from your local files"] = "  syncedpz add = adiciona um novo servidor PZ sincronizado a partir de seus arquivos locais"
	ptbrDict["  syncedpz delete = deletes a synced PZ server from the database only"] = "  syncedpz delete = exclui um servidor PZ sincronizado apenas do banco de dados"
	ptbrDict["  syncedpz clone = adds a new synced PZ server from a git repository"] = "  syncedpz clone = adiciona um novo servidor PZ sincronizado de um repositório git"
	ptbrDict["  syncedpz sync = syncs all servers"] = "  syncedpz sync = sincroniza todos os servidores"
	ptbrDict["  syncedpz play = syncs all servers at the start, every 5 minutes and at the end. And starts Project Zomboid"] = "  syncedpz play = sincroniza todo servidor no início, a cada 5 minutos e no final. E inicia o Project Zomboid"
	ptbrDict["  syncedpz language = sets the language of the application"] = "  syncedpz language = define o idioma da aplicação"
	ptbrDict["Menu:"] = "Menu:"
	ptbrDict["  [0] Help"] = "  [0] Ajuda"
	ptbrDict["  [1] Setup config"] = "  [1] Configuração"
	ptbrDict["  [2] List config"] = "  [2] Listar configuração"
	ptbrDict["  [3] List local servers"] = "  [3] Listar servidores locais"
	ptbrDict["  [4] List synced servers"] = "  [4] Listar servidores sincronizados"
	ptbrDict["  [5] Add synced server"] = "  [5] Adicionar servidor sincronizado"
	ptbrDict["  [6] Delete synced server"] = "  [6] Excluir servidor sincronizado"
	ptbrDict["  [7] Clone synced server"] = "  [7] Clonar servidor sincronizado"
	ptbrDict["  [8] Sync servers"] = "  [8] Sincronizar servidores"
	ptbrDict["  [9] Play"] = "  [9] Jogar"
	ptbrDict["  [10] Set language"] = "  [10] Definir idioma"
	ptbrDict["  [11] Exit"] = "  [11] Sair"
	ptbrDict["Enter the number of the option you want to choose: "] = "Digite o número da opção que deseja escolher: "
	ptbrDict["Invalid choice"] = "Escolha inválida"
	ptbrDict["Leave the field empty to use the previous value (if it exists)"] = "Deixe o campo vazio para usar o valor anterior (se existir)"
	ptbrDict["Enter the path to the pz executable (.bat file): "] = "Digite o caminho para o executável do pz (arquivo .bat): "
	ptbrDict["Enter the path to the pz data directory: "] = "Digite o caminho para a pasta de dados do pz: "
	ptbrDict["Enter your steam id: "] = "Digite seu id da steam: "
	ptbrDict["Enter your git username: "] = "Digite seu nome de usuário do git: "
	ptbrDict["Enter your git password (or your github token)): "] = "Digite sua senha do git (ou seu token do github): "
	ptbrDict["PZ Bat Path: "] = "Caminho do executável Bat executável Bat do PZ: "
	ptbrDict["PZ Data Path: "] = "Caminho dos dados do PZ: "
	ptbrDict["Steam ID: "] = "ID da Steam: "
	ptbrDict["Local Servers:"] = "Servidores Locais:"
	ptbrDict["Enter the number of the server you want to add: "] = "Digite o número do servidor que deseja adicionar: "
	ptbrDict["Enter the git repository link to the server: "] = "Digite o link do repositório git para o servidor: "
	ptbrDict["Warning! Apparently a server using this git repository already exists"] = "Atenção! Aparentemente um servidor usando este repositório git já existe"
	ptbrDict["and it already has some content."] = "e ele já possui algum conteúdo."
	ptbrDict["Do you want to continue copying your local content to it?"] = "Você deseja continuar copiando seu conteúdo local para ele?"
	ptbrDict["Enter y/N: "] = "Digite y/N (y para sim, n para não): "
	ptbrDict["Aborting."] = "Abortando."
	ptbrDict["Server added successfully"] = "Servidor adicionado com sucesso"
	ptbrDict["No servers to delete"] = "Nenhum servidor para excluir"
	ptbrDict["Enter the number of the server you want to delete: "] = "Digite o número do servidor que deseja excluir: "
	ptbrDict["Enter the git repository link to the server: "] = "Digite o link do repositório git para o servidor: "
	ptbrDict["Server cloned successfully"] = "Servidor clonado com sucesso"
	ptbrDict["Enter the number of the language you want to choose: "] = "Digite o número do idioma que deseja escolher: "
	ptbrDict["WARNING: Commiting and pushing can take a while, please wait..."] = "AVISO: Comitar e fazer push pode demorar um pouco, por favor aguarde..."

	dict[LANG_PTBR] = ptbrDict
}
