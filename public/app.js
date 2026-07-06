const sampleEvents = [
    {
        name: "THE ROCK BANDS LIVE I VEGA",
        description: "Oplev det legendariske rockband live i Kbenhavn, nar de gaester Store VEGA pa deres europaiske tour.",
        image: "https://example.com/images/concert-poster.jpg",
        startDate: "2026-10-15T20:00:00+02:00",
        endDate: "2026-10-15T23:00:00+02:00",
        doorTime: "2026-10-15T19:00:00+02:00",
        location: {
            name: "VEGA (Store VEGA)",
            address: { streetAddress: "Enghavevej 40", addressLocality: "Kbh V", postalCode: "1674", addressCountry: "DK" },
            sameAs: "https://www.vega.dk"
        },
        venueType: "vega",
        performer: { name: "The Rock Bands", sameAs: "https://example.com/artists/the-rock-bands" },
        offers: { price: "350", priceCurrency: "DKK", availability: "https://schema.org/InStock", url: "https://example.com/tickets/rock-bands-vega" },
        organizer: { name: "Live Nation Denmark", url: "https://www.livenation.dk" }
    },
    {
        name: "JAZZ NIGHT AT COPENHAGEN",
        description: "En aften med de bedste jazzmusikere i byen.",
        image: "https://example.com/images/jazz-night.jpg",
        startDate: "2026-10-20T19:30:00+02:00",
        endDate: "2026-10-20T22:30:00+02:00",
        doorTime: "2026-10-20T19:00:00+02:00",
        location: {
            name: "Jazzhus Montmartre",
            address: { streetAddress: "Store Regnegade 19A", addressLocality: "Kobenhavn", postalCode: "1110", addressCountry: "DK" },
            sameAs: "https://jazzhusmontmartre.dk"
        },
        venueType: "jazz",
        performer: { name: "Various Artists", sameAs: "https://example.com/artists/jazz-collective" },
        offers: { price: "200", priceCurrency: "DKK", availability: "https://schema.org/InStock", url: "https://example.com/tickets/jazz-night" },
        organizer: { name: "Copenhagen Jazz Society", url: "https://www.copenhagenjazz.dk" }
    },
    {
        name: "ELECTRONIC BEATS FESTIVAL",
        description: "Danmarks storste elektroniske musikfestival.",
        image: "https://example.com/images/electronic-festival.jpg",
        startDate: "2026-11-05T22:00:00+01:00",
        endDate: "2026-11-06T06:00:00+01:00",
        doorTime: "2026-11-05T21:00:00+01:00",
        location: {
            name: "Refshaleoen",
            address: { streetAddress: "Refshalevej 167", addressLocality: "Kobenhavn", postalCode: "1432", addressCountry: "DK" },
            sameAs: "https://refshaleoen.dk"
        },
        venueType: "other",
        performer: { name: "International DJ Lineup", sameAs: "https://example.com/artists/electronic-lineup" },
        offers: { price: "450", priceCurrency: "DKK", availability: "https://schema.org/InStock", url: "https://example.com/tickets/electronic-beats" },
        organizer: { name: "Electronic Events DK", url: "https://www.electronicevents.dk" }
    },
    {
        name: "CLASSICAL SYMPHONY ORCHESTRA",
        description: "Kobenhavns Filharmoniske Orkester praesenterer en aften med Beethovens symfonier.",
        image: "https://example.com/images/classical-orchestra.jpg",
        startDate: "2026-11-12T19:30:00+01:00",
        endDate: "2026-11-12T22:00:00+01:00",
        doorTime: "2026-11-12T19:00:00+01:00",
        location: {
            name: "DR Koncerthuset",
            address: { streetAddress: "Ingerslevsgade 85", addressLocality: "Kobenhavn", postalCode: "1799", addressCountry: "DK" },
            sameAs: "https://drkoncerthuset.dk"
        },
        venueType: "drkon",
        performer: { name: "Copenhagen Philharmonic Orchestra", sameAs: "https://example.com/artists/copenhagen-philharmonic" },
        offers: { price: "280", priceCurrency: "DKK", availability: "https://schema.org/InStock", url: "https://example.com/tickets/classical-symphony" },
        organizer: { name: "DR SymfoniOrkestret", url: "https://www.dr.dk/symfoni" }
    }
];

let currentFilter = 'all';

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('da-DK', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' });
}

function formatPrice(price, currency) {
    return new Intl.NumberFormat('da-DK', { style: 'currency', currency: currency }).format(price);
}

function getAvailabilityStatus(availability) {
    if (availability === 'https://schema.org/InStock') return 'available';
    else if (availability === 'https://schema.org/SoldOut') return 'sold-out';
    return 'available';
}

function createEventCard(event) {
    const locationIcon = '<svg class="inline w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd"></path></svg>';
    const status = getAvailabilityStatus(event.offers.availability);
    const statusText = status === 'available' ? 'TICKETS AVAILABLE' : 'SOLD OUT';
    return '<article class="event-card">' +
        '<div class="event-date">DATE: ' + formatDate(event.startDate) + '</div>' +
        '<h2 class="event-title">' + event.name + '</h2>' +
        '<div class="event-venue">' + locationIcon + event.location.name + ' - ' + event.location.address.addressLocality + '</div>' +
        '<p class="event-description">' + event.description + '</p>' +
        '<div class="flex justify-between items-center mt-4">' +
        '<span class="text-cyber-cyan font-medium">PERFORMER: ' + event.performer.name + '</span>' +
        '<span class="event-price">PRICE: ' + formatPrice(event.offers.price, event.offers.priceCurrency) + '</span>' +
        '</div>' +
        '<div class="mt-4">' +
        '<a href="' + event.offers.url + '" target="_blank" class="pixel-button">GET TICKETS</a>' +
        '<span class="event-status ' + status + '">' + statusText + '</span>' +
        '</div>' +
        '<div class="promoter-info">PROMOTED BY: ' + event.organizer.name + '</div>' +
        '</article>';
}

function filterEvents(events, filter) {
    if (filter === 'all') return events;
    return events.filter(event => event.venueType === filter);
}

function renderEvents(events) {
    const container = document.getElementById('events-grid');
    container.innerHTML = events.map(createEventCard).join('');
}

function showEvents() {
    document.getElementById('skeleton-loader').classList.add('hidden');
    document.getElementById('events-container').classList.remove('hidden');
}

function initFilters() {
    const filterChips = document.querySelectorAll('.filter-chip');
    filterChips.forEach(chip => {
        chip.addEventListener('click', () => {
            filterChips.forEach(c => c.classList.remove('active'));
            chip.classList.add('active');
            currentFilter = chip.dataset.filter;
            const filteredEvents = filterEvents(sampleEvents, currentFilter);
            renderEvents(filteredEvents);
        });
    });
}

function init() {
    setTimeout(() => {
        const filteredEvents = filterEvents(sampleEvents, currentFilter);
        renderEvents(filteredEvents);
        showEvents();
        initFilters();
    }, 1500);
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    init();
}