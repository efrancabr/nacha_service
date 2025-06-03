// NACHA Web Application JavaScript

document.addEventListener('DOMContentLoaded', function() {
    console.log('🌐 NACHA Web Application loaded');
    
    // Initialize components
    initializeFormValidation();
    initializeFileUpload();
    initializeNavigation();
    initializeTooltips();
});

// Form validation
function initializeFormValidation() {
    const forms = document.querySelectorAll('form');
    
    forms.forEach(form => {
        form.addEventListener('submit', function(e) {
            if (!validateForm(form)) {
                e.preventDefault();
                showNotification('Please fill in all required fields correctly.', 'error');
            }
        });
    });
}

function validateForm(form) {
    const requiredFields = form.querySelectorAll('[required]');
    let isValid = true;
    
    requiredFields.forEach(field => {
        if (!field.value.trim()) {
            field.style.borderColor = 'var(--error-color)';
            isValid = false;
        } else {
            field.style.borderColor = 'var(--border-color)';
        }
        
        // Specific validations
        if (field.type === 'number' && field.value < 0) {
            field.style.borderColor = 'var(--error-color)';
            isValid = false;
        }
        
        // Routing number validation (9 digits)
        if (field.name.includes('routing') || field.name.includes('destination') || field.name.includes('origin')) {
            if (!/^\d{9}$/.test(field.value)) {
                field.style.borderColor = 'var(--error-color)';
                isValid = false;
            }
        }
    });
    
    return isValid;
}

// File upload enhancements
function initializeFileUpload() {
    const fileInputs = document.querySelectorAll('input[type="file"]');
    
    fileInputs.forEach(input => {
        input.addEventListener('change', function(e) {
            const file = e.target.files[0];
            if (file) {
                validateFileType(file);
                updateFileUploadUI(file);
            }
        });
    });
}

function validateFileType(file) {
    const validTypes = ['.txt', '.ach'];
    const fileName = file.name.toLowerCase();
    const isValid = validTypes.some(type => fileName.endsWith(type));
    
    if (!isValid) {
        showNotification('Please select a valid NACHA file (.txt or .ach)', 'error');
        return false;
    }
    
    // Check file size (max 10MB)
    if (file.size > 10 * 1024 * 1024) {
        showNotification('File size must be less than 10MB', 'error');
        return false;
    }
    
    return true;
}

function updateFileUploadUI(file) {
    const uploadArea = document.querySelector('.file-upload');
    if (uploadArea) {
        uploadArea.style.borderColor = 'var(--success-color)';
        uploadArea.style.background = 'rgba(16, 185, 129, 0.05)';
        
        const icon = uploadArea.querySelector('.file-upload-icon');
        if (icon) icon.textContent = '✅';
        
        const text = uploadArea.querySelector('h3');
        if (text) text.textContent = `Selected: ${file.name}`;
    }
}

// Navigation enhancements
function initializeNavigation() {
    const currentPath = window.location.pathname;
    const navLinks = document.querySelectorAll('.nav-link');
    
    navLinks.forEach(link => {
        if (link.getAttribute('href') === currentPath) {
            link.style.background = 'rgba(255, 255, 255, 0.2)';
        }
        
        link.addEventListener('click', function(e) {
            // Add loading indicator for navigation
            if (!e.ctrlKey && !e.metaKey) {
                showLoading();
            }
        });
    });
}

// Notification system
function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `message message-${type} fade-in`;
    notification.innerHTML = `<strong>${type === 'error' ? 'Error:' : 'Info:'}</strong> ${message}`;
    
    // Remove existing notifications
    const existing = document.querySelector('.message');
    if (existing) existing.remove();
    
    // Insert at top of main content
    const main = document.querySelector('main');
    if (main) {
        main.insertBefore(notification, main.firstChild);
        
        // Auto-hide after 5 seconds
        setTimeout(() => {
            if (notification.parentNode) {
                notification.remove();
            }
        }, 5000);
    }
}

// Loading indicator
function showLoading() {
    const loader = document.createElement('div');
    loader.id = 'global-loader';
    loader.innerHTML = `
        <div style="
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(255, 255, 255, 0.8);
            display: flex;
            justify-content: center;
            align-items: center;
            z-index: 9999;
        ">
            <div class="loading" style="width: 40px; height: 40px;"></div>
        </div>
    `;
    
    document.body.appendChild(loader);
    
    // Remove loader if page doesn't change within 5 seconds
    setTimeout(() => {
        const existingLoader = document.getElementById('global-loader');
        if (existingLoader) existingLoader.remove();
    }, 5000);
}

// Tooltip system
function initializeTooltips() {
    const elements = document.querySelectorAll('[data-tooltip]');
    
    elements.forEach(element => {
        element.addEventListener('mouseenter', showTooltip);
        element.addEventListener('mouseleave', hideTooltip);
    });
}

function showTooltip(e) {
    const text = e.target.getAttribute('data-tooltip');
    const tooltip = document.createElement('div');
    tooltip.className = 'tooltip';
    tooltip.textContent = text;
    tooltip.style.cssText = `
        position: absolute;
        background: var(--text-primary);
        color: white;
        padding: 0.5rem;
        border-radius: var(--border-radius);
        font-size: 0.8rem;
        z-index: 1000;
        max-width: 200px;
        opacity: 0;
        transition: opacity 0.2s;
    `;
    
    document.body.appendChild(tooltip);
    
    // Position tooltip
    const rect = e.target.getBoundingClientRect();
    tooltip.style.left = rect.left + 'px';
    tooltip.style.top = (rect.bottom + 5) + 'px';
    
    // Fade in
    setTimeout(() => tooltip.style.opacity = '1', 10);
}

function hideTooltip() {
    const tooltip = document.querySelector('.tooltip');
    if (tooltip) tooltip.remove();
}

// Format helpers
function formatAmount(cents) {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD'
    }).format(cents / 100);
}

function formatDate(dateString) {
    if (dateString && dateString.length === 6) {
        const year = '20' + dateString.substring(0, 2);
        const month = dateString.substring(2, 4);
        const day = dateString.substring(4, 6);
        return new Date(year, month - 1, day).toLocaleDateString();
    }
    return dateString;
}

// Export functionality enhancements
function enhanceExportPage() {
    const formatButtons = document.querySelectorAll('.format-btn');
    
    formatButtons.forEach(button => {
        button.addEventListener('click', function() {
            // Add visual feedback
            formatButtons.forEach(btn => btn.classList.remove('selected'));
            this.classList.add('selected');
            
            // Add loading state
            this.innerHTML = `<div class="loading"></div> Processing...`;
            
            setTimeout(() => {
                // Restore button content (this would be replaced by actual export)
                this.innerHTML = this.getAttribute('data-original-content') || 'Export';
            }, 2000);
        });
        
        // Store original content
        button.setAttribute('data-original-content', button.innerHTML);
    });
}

// Auto-save functionality for forms
function initializeAutoSave() {
    const forms = document.querySelectorAll('form');
    
    forms.forEach(form => {
        const inputs = form.querySelectorAll('input, textarea, select');
        
        inputs.forEach(input => {
            // Load saved value
            const savedValue = localStorage.getItem(`form_${input.name}`);
            if (savedValue && !input.value) {
                input.value = savedValue;
            }
            
            // Save on change
            input.addEventListener('input', function() {
                localStorage.setItem(`form_${this.name}`, this.value);
            });
        });
        
        // Clear saved data on successful submit
        form.addEventListener('submit', function() {
            inputs.forEach(input => {
                localStorage.removeItem(`form_${input.name}`);
            });
        });
    });
}

// Keyboard shortcuts
document.addEventListener('keydown', function(e) {
    // Ctrl/Cmd + K for search/navigation
    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        // Focus first input or navigation
        const firstInput = document.querySelector('input, textarea');
        if (firstInput) firstInput.focus();
    }
    
    // Escape to close modals/notifications
    if (e.key === 'Escape') {
        const notifications = document.querySelectorAll('.message');
        notifications.forEach(n => n.remove());
    }
});

// Initialize additional features when needed
if (document.querySelector('.export-formats')) {
    enhanceExportPage();
}

if (document.querySelector('form')) {
    initializeAutoSave();
}

// Error handling
window.addEventListener('error', function(e) {
    console.error('Application error:', e.error);
    showNotification('An unexpected error occurred. Please try again.', 'error');
});

// Performance monitoring
window.addEventListener('load', function() {
    const loadTime = performance.now();
    console.log(`⚡ Page loaded in ${Math.round(loadTime)}ms`);
}); 