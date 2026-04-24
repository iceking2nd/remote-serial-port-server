export default {
  menu: {
    serial: {
      port: "Port",
      baudRate: "Débit en bauds",
      dataBits: "Bits de données",
      stopBits: "Bits d'arrêt",
      parity: "Parité",
      flowControl: "Contrôle de flux",
      placeholders: {
        port: "Veuillez sélectionner un port",
        baudRate: "Veuillez entrer le débit en bauds",
        dataBits: "Veuillez sélectionner les bits de données",
        stopBits: "Veuillez sélectionner les bits d'arrêt",
        parity: "Veuillez sélectionner la parité",
        flowControl: "Veuillez sélectionner le contrôle de flux",
      },
      options: {
        parity: {
          none: "Aucune",
          odd: "Impaire",
          even: "Paire",
          mark: "Marque",
          space: "Espace",
        }
      }
    },
    theme: {
      font: "Police",
      fontSize: "Taille de police",
      placeholders: {
        font: "Veuillez sélectionner une police",
        fontSize: "Veuillez entrer la taille de police",
      },
      buttons: {
        connect: "Connexion",
        disconnect: "Déconnexion",
      }
    }
  },
  message: {
    no_port: "Veuillez sélectionner un port",
    websocket_disconnected: "La connexion WebSocket est interrompue. Vérifiez le réseau ou reconnectez-vous.",
    websocket_error: "Une erreur WebSocket s'est produite. Vérifiez le réseau ou reconnectez-vous.",
    websocket_params_changed: "Les paramètres du formulaire ont changé, la connexion WebSocket est interrompue."
  },
  log: {
    api_key_not_available: "La clé API n'est pas disponible. Veuillez d'abord entrer la clé API.",
    websocket_connected: "WebSocket est connecté.",
    websocket_disconnected: "WebSocket est déconnecté.",
    websocket_error: "Erreur WebSocket :",
    websocket_manually_disconnect: "WebSocket a été déconnecté manuellement.",
  }
}
