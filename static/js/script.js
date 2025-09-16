// Modern Odoo Signup Form Handler
class OdooSignup {
    constructor() {
        this.form = document.getElementById('signupForm');
        this.usernameInput = document.getElementById('username');
        this.previewUrl = document.getElementById('preview-url');
        this.submitBtn = document.getElementById('submitBtn');
        this.loadingModal = document.getElementById('loadingModal');
        this.successModal = document.getElementById('successModal');
        this.progressFill = document.getElementById('progressFill');

        // Get domain from the suffix element
        this.domain = document.querySelector('.suffix').textContent.replace('.', '');

        this.init();
    }

    init() {
        this.setupEventListeners();
        this.setupFormValidation();
        this.updateUrlPreview();
        this.populateCountries();
        this.setFooterYear();
    }

    setFooterYear() {
        document.getElementById('currentYear').textContent = new Date().getFullYear();
    }

    setupEventListeners() {
        // Real-time URL preview
        this.usernameInput.addEventListener('input', () => {
            this.updateUrlPreview();
            this.validateUsername();
        });

        // Form submission
        this.form.addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleSubmit();
        });

        // Input validation on blur
        const inputs = this.form.querySelectorAll('input, select');
        inputs.forEach(input => {
            input.addEventListener('blur', () => this.validateField(input));
            input.addEventListener('input', () => this.clearFieldError(input));
        });
    }

    setupFormValidation() {
        // Custom validation rules
        this.validators = {
            username: (value) => {
                if (!value) return 'Username is required';
                if (value.length < 3) return 'Username must be at least 3 characters';
                if (value.length > 20) return 'Username must be less than 20 characters';
                if (!/^[a-zA-Z0-9_-]+$/.test(value)) return 'Username can only contain letters, numbers, hyphens, and underscores';
                if (value.toLowerCase() === 'admin' || value.toLowerCase() === 'www') return 'This username is not allowed';
                return null;
            },
            email: (value) => {
                if (!value) return 'Email is required';
                const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
                if (!emailRegex.test(value)) return 'Please enter a valid email address';
                return null;
            },
            password: (value) => {
                if (!value) return 'Password is required';
                if (value.length < 8) return 'Password must be at least 8 characters';
                if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(value)) {
                    return 'Password must contain at least one uppercase letter, one lowercase letter, and one number';
                }
                return null;
            },
            firstName: (value) => {
                if (!value) return 'First name is required';
                if (value.length < 2) return 'First name must be at least 2 characters';
                return null;
            },
            lastName: (value) => {
                if (!value) return 'Last name is required';
                if (value.length < 2) return 'Last name must be at least 2 characters';
                return null;
            },
            companyName: (value) => {
                if (!value) return 'Company name is required';
                if (value.length < 2) return 'Company name must be at least 2 characters';
                return null;
            },
            country: (value) => {
                if (!value) return 'Country is required';
                return null;
            },
            terms: (value) => {
                if (!value) return 'You must agree to the terms and conditions';
                return null;
            }
        };
    }

    populateCountries() {
        const countries = [
            {id: 15, name: "Åland Islands", code: "AX"},
            {id: 3, name: "Afghanistan", code: "AF"},
            {id: 6, name: "Albania", code: "AL"},
            {id: 62, name: "Algeria", code: "DZ"},
            {id: 11, name: "American Samoa", code: "AS"},
            {id: 1, name: "Andorra", code: "AD"},
            {id: 8, name: "Angola", code: "AO"},
            {id: 5, name: "Anguilla", code: "AI"},
            {id: 9, name: "Antarctica", code: "AQ"},
            {id: 4, name: "Antigua and Barbuda", code: "AG"},
            {id: 10, name: "Argentina", code: "AR"},
            {id: 7, name: "Armenia", code: "AM"},
            {id: 14, name: "Aruba", code: "AW"},
            {id: 13, name: "Australia", code: "AU"},
            {id: 12, name: "Austria", code: "AT"},
            {id: 16, name: "Azerbaijan", code: "AZ"},
            {id: 32, name: "Bahamas", code: "BS"},
            {id: 23, name: "Bahrain", code: "BH"},
            {id: 19, name: "Bangladesh", code: "BD"},
            {id: 18, name: "Barbados", code: "BB"},
            {id: 36, name: "Belarus", code: "BY"},
            {id: 20, name: "Belgium", code: "BE"},
            {id: 37, name: "Belize", code: "BZ"},
            {id: 25, name: "Benin", code: "BJ"},
            {id: 27, name: "Bermuda", code: "BM"},
            {id: 33, name: "Bhutan", code: "BT"},
            {id: 29, name: "Bolivia", code: "BO"},
            {id: 30, name: "Bonaire, Sint Eustatius and Saba", code: "BQ"},
            {id: 17, name: "Bosnia and Herzegovina", code: "BA"},
            {id: 35, name: "Botswana", code: "BW"},
            {id: 34, name: "Bouvet Island", code: "BV"},
            {id: 31, name: "Brazil", code: "BR"},
            {id: 105, name: "British Indian Ocean Territory", code: "IO"},
            {id: 28, name: "Brunei Darussalam", code: "BN"},
            {id: 22, name: "Bulgaria", code: "BG"},
            {id: 21, name: "Burkina Faso", code: "BF"},
            {id: 24, name: "Burundi", code: "BI"},
            {id: 116, name: "Cambodia", code: "KH"},
            {id: 47, name: "Cameroon", code: "CM"},
            {id: 38, name: "Canada", code: "CA"},
            {id: 52, name: "Cape Verde", code: "CV"},
            {id: 123, name: "Cayman Islands", code: "KY"},
            {id: 40, name: "Central African Republic", code: "CF"},
            {id: 214, name: "Chad", code: "TD"},
            {id: 46, name: "Chile", code: "CL"},
            {id: 48, name: "China", code: "CN"},
            {id: 54, name: "Christmas Island", code: "CX"},
            {id: 39, name: "Cocos (Keeling) Islands", code: "CC"},
            {id: 49, name: "Colombia", code: "CO"},
            {id: 118, name: "Comoros", code: "KM"},
            {id: 42, name: "Congo", code: "CG"},
            {id: 45, name: "Cook Islands", code: "CK"},
            {id: 50, name: "Costa Rica", code: "CR"},
            {id: 97, name: "Croatia", code: "HR"},
            {id: 51, name: "Cuba", code: "CU"},
            {id: 53, name: "Curaçao", code: "CW"},
            {id: 55, name: "Cyprus", code: "CY"},
            {id: 56, name: "Czech Republic", code: "CZ"},
            {id: 44, name: "Côte d'Ivoire", code: "CI"},
            {id: 41, name: "Democratic Republic of the Congo", code: "CD"},
            {id: 59, name: "Denmark", code: "DK"},
            {id: 58, name: "Djibouti", code: "DJ"},
            {id: 60, name: "Dominica", code: "DM"},
            {id: 61, name: "Dominican Republic", code: "DO"},
            {id: 63, name: "Ecuador", code: "EC"},
            {id: 65, name: "Egypt", code: "EG"},
            {id: 209, name: "El Salvador", code: "SV"},
            {id: 87, name: "Equatorial Guinea", code: "GQ"},
            {id: 67, name: "Eritrea", code: "ER"},
            {id: 64, name: "Estonia", code: "EE"},
            {id: 212, name: "Eswatini", code: "SZ"},
            {id: 69, name: "Ethiopia", code: "ET"},
            {id: 72, name: "Falkland Islands", code: "FK"},
            {id: 74, name: "Faroe Islands", code: "FO"},
            {id: 71, name: "Fiji", code: "FJ"},
            {id: 70, name: "Finland", code: "FI"},
            {id: 75, name: "France", code: "FR"},
            {id: 79, name: "French Guiana", code: "GF"},
            {id: 174, name: "French Polynesia", code: "PF"},
            {id: 215, name: "French Southern Territories", code: "TF"},
            {id: 76, name: "Gabon", code: "GA"},
            {id: 84, name: "Gambia", code: "GM"},
            {id: 78, name: "Georgia", code: "GE"},
            {id: 57, name: "Germany", code: "DE"},
            {id: 80, name: "Ghana", code: "GH"},
            {id: 81, name: "Gibraltar", code: "GI"},
            {id: 88, name: "Greece", code: "GR"},
            {id: 83, name: "Greenland", code: "GL"},
            {id: 77, name: "Grenada", code: "GD"},
            {id: 86, name: "Guadeloupe", code: "GP"},
            {id: 91, name: "Guam", code: "GU"},
            {id: 90, name: "Guatemala", code: "GT"},
            {id: 82, name: "Guernsey", code: "GG"},
            {id: 85, name: "Guinea", code: "GN"},
            {id: 92, name: "Guinea-Bissau", code: "GW"},
            {id: 93, name: "Guyana", code: "GY"},
            {id: 98, name: "Haiti", code: "HT"},
            {id: 95, name: "Heard Island and McDonald Islands", code: "HM"},
            {id: 236, name: "Holy See (Vatican City State)", code: "VA"},
            {id: 96, name: "Honduras", code: "HN"},
            {id: 94, name: "Hong Kong", code: "HK"},
            {id: 99, name: "Hungary", code: "HU"},
            {id: 108, name: "Iceland", code: "IS"},
            {id: 104, name: "India", code: "IN"},
            {id: 100, name: "Indonesia", code: "ID"},
            {id: 107, name: "Iran", code: "IR"},
            {id: 106, name: "Iraq", code: "IQ"},
            {id: 101, name: "Ireland", code: "IE"},
            {id: 103, name: "Isle of Man", code: "IM"},
            {id: 102, name: "Israel", code: "IL"},
            {id: 109, name: "Italy", code: "IT"},
            {id: 111, name: "Jamaica", code: "JM"},
            {id: 113, name: "Japan", code: "JP"},
            {id: 110, name: "Jersey", code: "JE"},
            {id: 112, name: "Jordan", code: "JO"},
            {id: 124, name: "Kazakhstan", code: "KZ"},
            {id: 114, name: "Kenya", code: "KE"},
            {id: 117, name: "Kiribati", code: "KI"},
            {id: 250, name: "Kosovo", code: "XK"},
            {id: 122, name: "Kuwait", code: "KW"},
            {id: 115, name: "Kyrgyzstan", code: "KG"},
            {id: 125, name: "Laos", code: "LA"},
            {id: 134, name: "Latvia", code: "LV"},
            {id: 126, name: "Lebanon", code: "LB"},
            {id: 131, name: "Lesotho", code: "LS"},
            {id: 130, name: "Liberia", code: "LR"},
            {id: 135, name: "Libya", code: "LY"},
            {id: 128, name: "Liechtenstein", code: "LI"},
            {id: 132, name: "Lithuania", code: "LT"},
            {id: 133, name: "Luxembourg", code: "LU"},
            {id: 147, name: "Macau", code: "MO"},
            {id: 141, name: "Madagascar", code: "MG"},
            {id: 155, name: "Malawi", code: "MW"},
            {id: 157, name: "Malaysia", code: "MY"},
            {id: 154, name: "Maldives", code: "MV"},
            {id: 144, name: "Mali", code: "ML"},
            {id: 152, name: "Malta", code: "MT"},
            {id: 142, name: "Marshall Islands", code: "MH"},
            {id: 149, name: "Martinique", code: "MQ"},
            {id: 150, name: "Mauritania", code: "MR"},
            {id: 153, name: "Mauritius", code: "MU"},
            {id: 246, name: "Mayotte", code: "YT"},
            {id: 156, name: "Mexico", code: "MX"},
            {id: 73, name: "Micronesia", code: "FM"},
            {id: 138, name: "Moldova", code: "MD"},
            {id: 137, name: "Monaco", code: "MC"},
            {id: 146, name: "Mongolia", code: "MN"},
            {id: 139, name: "Montenegro", code: "ME"},
            {id: 151, name: "Montserrat", code: "MS"},
            {id: 136, name: "Morocco", code: "MA"},
            {id: 158, name: "Mozambique", code: "MZ"},
            {id: 145, name: "Myanmar", code: "MM"},
            {id: 159, name: "Namibia", code: "NA"},
            {id: 168, name: "Nauru", code: "NR"},
            {id: 167, name: "Nepal", code: "NP"},
            {id: 165, name: "Netherlands", code: "NL"},
            {id: 160, name: "New Caledonia", code: "NC"},
            {id: 170, name: "New Zealand", code: "NZ"},
            {id: 164, name: "Nicaragua", code: "NI"},
            {id: 161, name: "Niger", code: "NE"},
            {id: 163, name: "Nigeria", code: "NG"},
            {id: 169, name: "Niue", code: "NU"},
            {id: 162, name: "Norfolk Island", code: "NF"},
            {id: 120, name: "North Korea", code: "KP"},
            {id: 143, name: "North Macedonia", code: "MK"},
            {id: 148, name: "Northern Mariana Islands", code: "MP"},
            {id: 166, name: "Norway", code: "NO"},
            {id: 171, name: "Oman", code: "OM"},
            {id: 177, name: "Pakistan", code: "PK"},
            {id: 184, name: "Palau", code: "PW"},
            {id: 172, name: "Panama", code: "PA"},
            {id: 175, name: "Papua New Guinea", code: "PG"},
            {id: 185, name: "Paraguay", code: "PY"},
            {id: 173, name: "Peru", code: "PE"},
            {id: 176, name: "Philippines", code: "PH"},
            {id: 180, name: "Pitcairn Islands", code: "PN"},
            {id: 178, name: "Poland", code: "PL"},
            {id: 183, name: "Portugal", code: "PT"},
            {id: 181, name: "Puerto Rico", code: "PR"},
            {id: 186, name: "Qatar", code: "QA"},
            {id: 188, name: "Romania", code: "RO"},
            {id: 190, name: "Russian Federation", code: "RU"},
            {id: 191, name: "Rwanda", code: "RW"},
            {id: 187, name: "Réunion", code: "RE"},
            {id: 26, name: "Saint Barthélémy", code: "BL"},
            {id: 198, name: "Saint Helena, Ascension and Tristan da Cunha", code: "SH"},
            {id: 119, name: "Saint Kitts and Nevis", code: "KN"},
            {id: 127, name: "Saint Lucia", code: "LC"},
            {id: 140, name: "Saint Martin (French part)", code: "MF"},
            {id: 179, name: "Saint Pierre and Miquelon", code: "PM"},
            {id: 237, name: "Saint Vincent and the Grenadines", code: "VC"},
            {id: 244, name: "Samoa", code: "WS"},
            {id: 203, name: "San Marino", code: "SM"},
            {id: 192, name: "Saudi Arabia", code: "SA"},
            {id: 204, name: "Senegal", code: "SN"},
            {id: 189, name: "Serbia", code: "RS"},
            {id: 194, name: "Seychelles", code: "SC"},
            {id: 202, name: "Sierra Leone", code: "SL"},
            {id: 197, name: "Singapore", code: "SG"},
            {id: 210, name: "Sint Maarten (Dutch part)", code: "SX"},
            {id: 201, name: "Slovakia", code: "SK"},
            {id: 199, name: "Slovenia", code: "SI"},
            {id: 193, name: "Solomon Islands", code: "SB"},
            {id: 205, name: "Somalia", code: "SO"},
            {id: 247, name: "South Africa", code: "ZA"},
            {id: 89, name: "South Georgia and the South Sandwich Islands", code: "GS"},
            {id: 121, name: "South Korea", code: "KR"},
            {id: 207, name: "South Sudan", code: "SS"},
            {id: 68, name: "Spain", code: "ES"},
            {id: 129, name: "Sri Lanka", code: "LK"},
            {id: 182, name: "State of Palestine", code: "PS"},
            {id: 195, name: "Sudan", code: "SD"},
            {id: 206, name: "Suriname", code: "SR"},
            {id: 200, name: "Svalbard and Jan Mayen", code: "SJ"},
            {id: 196, name: "Sweden", code: "SE"},
            {id: 43, name: "Switzerland", code: "CH"},
            {id: 211, name: "Syria", code: "SY"},
            {id: 208, name: "São Tomé and Príncipe", code: "ST"},
            {id: 227, name: "Taiwan", code: "TW"},
            {id: 218, name: "Tajikistan", code: "TJ"},
            {id: 228, name: "Tanzania", code: "TZ"},
            {id: 217, name: "Thailand", code: "TH"},
            {id: 223, name: "Timor-Leste", code: "TL"},
            {id: 216, name: "Togo", code: "TG"},
            {id: 219, name: "Tokelau", code: "TK"},
            {id: 222, name: "Tonga", code: "TO"},
            {id: 225, name: "Trinidad and Tobago", code: "TT"},
            {id: 221, name: "Tunisia", code: "TN"},
            {id: 220, name: "Turkmenistan", code: "TM"},
            {id: 213, name: "Turks and Caicos Islands", code: "TC"},
            {id: 226, name: "Tuvalu", code: "TV"},
            {id: 224, name: "Türkiye", code: "TR"},
            {id: 232, name: "USA Minor Outlying Islands", code: "UM"},
            {id: 230, name: "Uganda", code: "UG"},
            {id: 229, name: "Ukraine", code: "UA"},
            {id: 2, name: "United Arab Emirates", code: "AE"},
            {id: 231, name: "United Kingdom", code: "GB"},
            {id: 233, name: "United States", code: "US"},
            {id: 234, name: "Uruguay", code: "UY"},
            {id: 235, name: "Uzbekistan", code: "UZ"},
            {id: 242, name: "Vanuatu", code: "VU"},
            {id: 238, name: "Venezuela", code: "VE"},
            {id: 241, name: "Vietnam", code: "VN"},
            {id: 239, name: "Virgin Islands (British)", code: "VG"},
            {id: 240, name: "Virgin Islands (USA)", code: "VI"},
            {id: 243, name: "Wallis and Futuna", code: "WF"},
            {id: 66, name: "Western Sahara", code: "EH"},
            {id: 245, name: "Yemen", code: "YE"},
            {id: 248, name: "Zambia", code: "ZM"},
            {id: 249, name: "Zimbabwe", code: "ZW"}
        ];

        const countrySelect = document.getElementById('country');
        countries.sort((a, b) => a.name.localeCompare(b.name)).forEach(country => {
            const option = document.createElement('option');
            option.dataset.id = country.id;
            option.value = country.code;
            option.textContent = country.name;
            
            // Set Philippines as default
            if (country.name === "Philippines") {
                option.selected = true;
            }
            
            countrySelect.appendChild(option);
        });
    }

    updateUrlPreview() {
        const username = this.usernameInput.value.trim();
        if (username) {
            this.previewUrl.textContent = `${username}.${this.domain}`;
        } else {
            this.previewUrl.textContent = `yourcompany.${this.domain}`;
        }
    }

    validateUsername() {
        const error = this.validators.username(this.usernameInput.value);
        this.showFieldError(this.usernameInput, error);
        return !error;
    }

    validateField(field) {
        const validator = this.validators[field.name];
        if (validator) {
            const error = validator(field.value);
            this.showFieldError(field, error);
            return !error;
        }
        return true;
    }

    showFieldError(field, error) {
        this.clearFieldError(field);
        if (error) {
            field.classList.add('error');
            const errorElement = document.createElement('div');
            errorElement.className = 'field-error';
            errorElement.innerHTML = `<i class="fas fa-exclamation-circle"></i> ${error}`;
            field.parentNode.appendChild(errorElement);
        }
    }

    clearFieldError(field) {
        field.classList.remove('error');
        const errorElement = field.parentNode.querySelector('.field-error');
        if (errorElement) {
            errorElement.remove();
        }
    }

    validateForm() {
        let isValid = true;
        const inputs = this.form.querySelectorAll('input, select');

        inputs.forEach(input => {
            if (!this.validateField(input)) {
                isValid = false;
            }
        });

        return isValid;
    }

    async handleSubmit() {
        if (!this.validateForm()) {
            this.showNotification('Please fix the errors in the form', 'error');
            return;
        }

        this.showLoadingModal();
        this.setLoadingState(true);

        try {
            const formData = this.getFormData();
            const response = await this.submitSignup(formData);

            if (response.success) {
                this.showSuccessModal(response.data);
                this.clearForm();
            } else {
                throw new Error(response.message || 'Signup failed');
            }
        } catch (error) {
            console.error('Signup error:', error);
            this.showNotification(error.message || 'An error occurred during signup', 'error');
        } finally {
            this.setLoadingState(false);
            this.hideLoadingModal();
        }
    }

    getFormData() {
        const formData = new FormData(this.form);
        const countrySelect = document.getElementById('country');
        const selectedOption = countrySelect.selectedOptions[0];
        return {
            username: formData.get('username').trim(),
            email: formData.get('email').trim(),
            password: formData.get('password'),
            firstName: formData.get('firstName').trim(),
            lastName: formData.get('lastName').trim(),
            phone: formData.get('phone').trim(),
            companyName: formData.get('companyName').trim(),
            industry: formData.get('industry'),
            companySize: formData.get('companySize'),
            country: selectedOption ? {
                id: parseInt(selectedOption.dataset.id),
                code: selectedOption.value,
                name: selectedOption.textContent
            } : null,
            terms: formData.get('terms') === 'on'
        };
    }

    async submitSignup(formData) {
        // Simulate progress updates
        this.updateProgress(20);

        const response = await fetch('/api/signup', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        });

        this.updateProgress(60);

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.message || `HTTP ${response.status}`);
        }

        this.updateProgress(100);
        return await response.json();
    }

    showLoadingModal() {
        this.loadingModal.style.display = 'block';
        document.body.style.overflow = 'hidden';
        this.updateProgress(0);
    }

    hideLoadingModal() {
        this.loadingModal.style.display = 'none';
        document.body.style.overflow = '';
    }

    showSuccessModal(data) {
        document.getElementById('instanceUrl').href = `https://${data.instanceUrl}`;
        document.getElementById('instanceUrl').textContent = data.instanceUrl;
        document.getElementById('adminEmail').textContent = data.email;

        // Set up the access instance button
        const accessBtn = document.getElementById('accessInstanceBtn');
        accessBtn.onclick = () => {
            window.open(`https://${data.instanceUrl}`, '_blank');
        };

        this.successModal.style.display = 'block';
        document.body.style.overflow = 'hidden';
    }

    clearForm() {
        this.form.reset();
        this.updateUrlPreview();
        // Clear any error messages
        const errorElements = this.form.querySelectorAll('.field-error');
        errorElements.forEach(element => element.remove());
        // Remove error classes from inputs
        const errorInputs = this.form.querySelectorAll('.error');
        errorInputs.forEach(input => input.classList.remove('error'));
    }

    updateProgress(percentage) {
        this.progressFill.style.width = `${percentage}%`;
    }

    setLoadingState(loading) {
        this.submitBtn.disabled = loading;
        this.submitBtn.innerHTML = loading
            ? '<i class="fas fa-spinner fa-spin"></i> Creating Instance...'
            : '<span class="btn-text">Create My Odoo Instance</span><i class="fas fa-arrow-right"></i>';

        if (loading) {
            this.form.classList.add('loading');
        } else {
            this.form.classList.remove('loading');
        }
    }

    showNotification(message, type = 'info') {
        // Remove existing notifications
        const existingNotifications = document.querySelectorAll('.notification');
        existingNotifications.forEach(notification => notification.remove());

        // Create new notification
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.innerHTML = `
            <i class="fas fa-${type === 'error' ? 'exclamation-circle' : 'check-circle'}"></i>
            ${message}
            <button class="notification-close" onclick="this.parentNode.remove()">
                <i class="fas fa-times"></i>
            </button>
        `;

        // Add to page
        document.body.appendChild(notification);

        // Auto remove after 5 seconds
        setTimeout(() => {
            if (notification.parentNode) {
                notification.remove();
            }
        }, 5000);
    }
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new OdooSignup();
});