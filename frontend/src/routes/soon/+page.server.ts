export async function load() {
    try {
        // const res = await fetch('http://backend:3500/hello');
        // TODO: does that thing work, guess it would be better to use this ?
        const domain = process.env.API_URL
        const res = await fetch(`http://${domain}/hello`);
        if (!res.ok) {
            throw new Error(`HTTP error status: ${res.status}`);
        }
        const message = await res.json();
        console.log('the message is:', message);
    } catch (err) {
        console.error('Error fetching data:', err);
    }
}
