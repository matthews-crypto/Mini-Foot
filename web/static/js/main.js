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

// Fonction pour faire une requête API
async function apiRequest(url, method, body = null) {
    const headers = {
        'Content-Type': 'application/json',
    };

    const token = getToken();
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const options = {
        method,
        headers,
    };

    if (body) {
        options.body = JSON.stringify(body);
    }

    const response = await fetch(url, options);
    const data = await response.json();

    if (!response.ok) {
        throw new Error(data.message || 'Une erreur est survenue');
    }

    return data;
}

// Gestionnaire d'inscription et connexion
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
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ nom, prenom, telephone, email, mot_de_passe: mot_de_passe }),
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        console.log('Réponse du serveur:', data);
        alert('Inscription réussie ! Vous pouvez maintenant vous connecter.');
        window.location.href = '/login';
    } catch (error) {
        console.error('Erreur lors de l\'inscription:', error);
        alert('Erreur lors de l\'inscription. Veuillez réessayer.');
    }
}

async function handleLogin(e) {
    e.preventDefault();
    
    const telephoneInput = document.getElementById('loginTelephone');
    const passwordInput = document.getElementById('loginPassword');

    if (!telephoneInput || !passwordInput) {
        console.error('Les champs de formulaire n\'ont pas été trouvés');
        alert('Erreur lors de la connexion. Veuillez réessayer.');
        return;
    }

    const telephone = telephoneInput.value;
    const mot_de_passe = passwordInput.value;

    console.log('Tentative de connexion avec:', { 
        telephone, 
        mot_de_passe: mot_de_passe ? `[REMPLI: ${mot_de_passe.length} caractères]` : '[VIDE]' 
    });

    if (!telephone || !mot_de_passe) {
        alert('Veuillez remplir tous les champs');
        return;
    }

    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ telephone, mot_de_passe }),
        });

        console.log('Corps de la requête:', JSON.stringify({ telephone, mot_de_passe }));

        if (!response.ok) {
            const errorText = await response.text();
            console.error('Réponse du serveur en cas d\'erreur:', errorText);
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        console.log('Réponse du serveur:', data);
        if (data.token) {
            localStorage.setItem('token', data.token);
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

// Assurez-vous que cet écouteur d'événements est bien en place
document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    } else {
        console.error('Le formulaire de connexion n\'a pas été trouvé');
    }
});

// N'oubliez pas d'ajouter cet écouteur d'événements à la fin du fichier
document.addEventListener('DOMContentLoaded', function() {
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', handleRegister);
    }
});

// N'oubliez pas d'ajouter cet écouteur d'événements à la fin du fichier
document.addEventListener('DOMContentLoaded', function() {
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', handleRegister);
    }
});

// Gestionnaire de connexion


// N'oubliez pas d'ajouter cet écouteur d'événements à la fin du fichier si ce n'est pas déjà fait
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
        await apiRequest('/profile', 'PUT', { nom, prenom, telephone, email });
        alert('Profil mis à jour avec succès !');
    } catch (error) {
        alert(error.message);
    }
}

// Gestionnaire de suppression de compte
async function handleDeleteAccount() {
    if (confirm('Êtes-vous sûr de vouloir supprimer votre compte ? Cette action est irréversible.')) {
        try {
            await apiRequest('/account', 'DELETE');
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
    try {
        const userData = await apiRequest('/profile', 'GET');
        document.getElementById('profileNom').value = userData.nom;
        document.getElementById('profilePrenom').value = userData.prenom;
        document.getElementById('profileTelephone').value = userData.telephone;
        document.getElementById('profileEmail').value = userData.email || '';
    } catch (error) {
        alert(error.message);
    }
}

// Ajouter les gestionnaires d'événements aux formulaires
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
        loadProfileInfo();
    }

    const deleteAccountButton = document.getElementById('deleteAccount');
    if (deleteAccountButton) {
        deleteAccountButton.addEventListener('click', handleDeleteAccount);
    }
});