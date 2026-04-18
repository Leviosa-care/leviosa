export const mockSettings = {
    company: {
        name: "Leviosa Spa & Wellness Center",
        email: "contact@leviosa-spa.com",
        telephone: "+33123456789",
        address: "123 Wellness Boulevard, 75001 Paris, France",
        instagram: "https://instagram.com/leviosa_spa",
        logo_url: "https://placehold.co/400x200/png?text=Leviosa+Logo",
        logo_content_type: "image/png"
    },
    otp: {
        duration: 300,        // seconds (5 minutes)
        length: 6,            // digits
        max_attempts: 5       // attempts
    },
    tokens: {
        access_duration: 15,  // minutes
        refresh_duration: 168 // hours (7 days)
    }
}
