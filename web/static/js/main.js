// Fonction pour obtenir le token JWT du localStorage
function getToken() {
    return localStorage.getItem('token');
}

// Fonction pour définir le token JWT dans le localStorage
function setToken(token) {
    localStorage.setItem('token', token);
}

// Fonction pour supprimer le token JWT du localStorage
function removeToken() {
    localStorage.removeItem('token');
}

// Fonction pour faire une requête API authentifiée
async function apiRequest(url, method, body = null) {
    const headers = {
        'Content-Type': 'application/json',
    };

    const token = localStorage.getItem('token');
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
        console.log('Token envoyé:', token);
    } else {
        console.log('Aucun token trouvé dans le localStorage');
    }

    const options = {
        method,
        headers,
    };

    if (body) {
        options.body = JSON.stringify(body);
    }

    console.log('Requête envoyée à:', url, 'avec options:', options);

    const response = await fetch(url, options);
    
    if (!response.ok) {
        const errorText = await response.text();
        console.error('Erreur de réponse:', response.status, errorText);
        throw new Error(errorText || 'Une erreur est survenue');
    }

    return response.json();
}

// Gestionnaire d'inscription
async function handleRegister(e) {
    e.preventDefault();
    const nom = document.getElementById('nom').value;
    const prenom = document.getElementById('prenom').value;
    const telephone = document.getElementById('telephone').value;
    const email = document.getElementById('email').value;
    const mot_de_passe = document.getElementById('mot_de_passe').value;

    console.log('Tentative d\'inscription avec:', { nom, prenom, telephone, email, mot_de_passe: mot_de_passe ? '[REMPLI]' : '[VIDE]' });

    if (!mot_de_passe) {
        alert('Le mot de passe ne peut pas être vide.');
        return;
    }

    try {
        const data = await apiRequest('/api/register', 'POST', { nom, prenom, telephone, email, mot_de_passe });
        console.log('Réponse du serveur:', data);
        alert('Inscription réussie ! Vous pouvez maintenant vous connecter.');
        window.location.href = '/login';
    } catch (error) {
        console.error('Erreur lors de l\'inscription:', error);
        alert('Erreur lors de l\'inscription. Veuillez réessayer.');
    }
}

// Gestionnaire de connexion
async function handleLogin(e) {
    e.preventDefault();
    
    const telephone = document.getElementById('loginTelephone').value;
    const mot_de_passe = document.getElementById('loginPassword').value;

    try {
        const data = await apiRequest('/api/login', 'POST', { telephone, mot_de_passe });
        console.log('Réponse du serveur:', data);
        if (data.token) {
            localStorage.setItem('token', data.token);
            console.log('Token stocké dans localStorage:', data.token);
            alert('Connexion réussie !');
            window.location.href = '/profile';
        } else {
            throw new Error('Token non reçu');
        }
    } catch (error) {
        console.error('Erreur lors de la connexion:', error);
        alert('Erreur lors de la connexion. Veuillez réessayer.');
    }
}

// Assurez-vous d'avoir cet écouteur d'événements
document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }
});

// Gestionnaire de mise à jour du profil
async function handleUpdateProfile(e) {
    e.preventDefault();
    const nom = document.getElementById('profileNom').value;
    const prenom = document.getElementById('profilePrenom').value;
    const telephone = document.getElementById('profileTelephone').value;
    const email = document.getElementById('profileEmail').value;

    try {
        await apiRequest('/api/profile', 'PUT', { nom, prenom, telephone, email });
        alert('Profil mis à jour avec succès !');
    } catch (error) {
        alert(error.message);
    }
}

// Gestionnaire de suppression de compte
async function handleDeleteAccount() {
    if (confirm('Êtes-vous sûr de vouloir supprimer votre compte ? Cette action est irréversible.')) {
        try {
            await apiRequest('/api/account', 'DELETE');
            removeToken();
            alert('Votre compte a été supprimé avec succès.');
            window.location.href = '/';
        } catch (error) {
            alert(error.message);
        }
    }
}

// Charger les informations du profil
async function loadProfileInfo() {
    console.log('Début de loadProfileInfo');
    try {
        const userData = await apiRequest('/api/profile', 'GET');
        console.log('Données du profil reçues:', userData);
        document.getElementById('profileNom').value = userData.nom;
        document.getElementById('profilePrenom').value = userData.prenom;
        document.getElementById('profileTelephone').value = userData.telephone;
        document.getElementById('profileEmail').value = userData.email || '';
    } catch (error) {
        console.error('Erreur détaillée lors du chargement du profil:', error);
        if (error.message.includes('401')) {
            console.log('Erreur d\'authentification détectée, redirection vers la page de connexion');
            window.location.href = '/login';
        } else {
            alert('Erreur lors du chargement du profil: ' + error.message);
        }
    }
}

// Initialisation des gestionnaires d'événements
document.addEventListener('DOMContentLoaded', function() {
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', handleRegister);
    }

    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }

    const profileForm = document.getElementById('profileForm');
    if (profileForm) {
        profileForm.addEventListener('submit', handleUpdateProfile);
    }

    const deleteAccountButton = document.getElementById('deleteAccount');
    if (deleteAccountButton) {
        deleteAccountButton.addEventListener('click', handleDeleteAccount);
    }

    // Charger le profil si on est sur la page de profil
    if (window.location.pathname === '/profile') {
        console.log('Page de profil détectée, chargement des informations du profil');
        loadProfileInfo();
    }
});