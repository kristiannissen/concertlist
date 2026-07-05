const sampleEvents = [
    {
        name: "The Rock Bands Live i Vega",
        description: "Oplev det legendariske rockband live i København, når de gæster Store VEGA på deres europæiske tour.",
        image: "https://example.com/images/concert-poster.jpg",
        startDate: "2026-10-15T20:00:00+02:00",
        endDate: "2026-10-15T23:00:00+02:00",
        doorTime: "2026-10-15T19:00:00+02:00",
        location: {
            name: "VEGA (Store VEGA)",
            address: {
                streetAddress: "Enghavevej 40",
                addressLocality: "København V",
                postalCode: "1674",
                addressCountry: "DK"
            },
            sameAs: "https://www.vega.dk"
        },
        performer: {
            name: "The Rock Bands",
            sameAs: "https://example.com/artists/the-rock-bands"
        },
        offers: {
            price: "350",
            priceCurrency: "DKK",
            availability: "https://schema.org/InStock",
            url: "https://example.com/tickets/rock-bands-vega"
        },
        organizer: {
            name: "Live Nation Denmark",
            url: "https://www.livenation.dk"
        }
    },
    {
        name: "Jazz Night at Copenhagen",
        description: "En aften med de bedste jazzmusikere i byen.",
        image: "https://example.com/images/jazz-night.jpg",
        startDate: "2026-10-20T19:30:00+02:00",
        endDate: "2026-10-20T22:30:00+02:00",
        doorTime: "2026-10-20T19:00:00+02:00",
        location: {
            name: "Jazzhus Montmartre",
            address: {
                streetAddress: "Store Regnegade 19A",
                addressLocality: "København",
                postalCode: "1110",
                addressCountry: "DK"
            },
            sameAs: "https://jazzhusmontmartre.dk"
        },
        performer: {
            name: "Various Artists",
            sameAs: "https://example.com/artists/jazz-collective"
        },
        offers: {
            price: "200",
            priceCurrency: "DKK",
            availability: "https://schema.org/InStock",
            url: "https://example.com/tickets/jazz-night"
        },
        organizer: {
            name: "Copenhagen Jazz Society",
            url: "https://www.copenhagenjazz.dk"
        }
    },
    {
        name: "Electronic Beats Festival",
        description: "Danmarks største elektroniske musikfestival.",
        image: "https://example.com/images/electronic-festival.jpg",
        startDate: "2026-11-05T22:00:00+01:00",
        endDate: "2026-11-06T06:00:00+01:00",
        doorTime: "2026-11-05T21:00:00+01:00",
        location: {
            name: "Refshaleøen",
            address: {
                streetAddress: "Refshalevej 167",
                addressLocality: "København",
                postalCode: "1432",
                addressCountry: "DK"
            },
            sameAs: "https://refshaleoen.dk"
        },
        performer: {
            name: "International DJ Lineup",
            sameAs: "https://example.com/artists/electronic-lineup"
        },
        offers: {
            price: "450",
            priceCurrency: "DKK",
            availability: "https://schema.org/InStock",
            url: "https://example.com/tickets/electronic-beats"
        },
        organizer: {
            name: "Electronic Events DK",
            url: "https://www.electronicevents.dk"
        }
    },
    {
        name: "Classical Symphony Orchestra",
        description: "Københavns Filharmoniske Orkester præsenterer en aften med Beethovens symfonier.",
        image: "https://example.com/images/classical-orchestra.jpg",
        startDate: "2026-11-12T19:30:00+01:00",
        endDate: "2026-11-12T22:00:00+01:00",
        doorTime: "2026-11-12T19:00:00+01:00",
        location: {
            name: "DR Koncerthuset",
            address: {
                streetAddress: "Ingerslevsgade 85",
                addressLocality: "København",
                postalCode: "1799",
                addressCountry: "DK"
            },
            sameAs: "https://drkoncerthuset.dk"
        },
        performer: {
            name: "Copenhagen Philharmonic Orchestra",
            sameAs: "https://example.com/artists/copenhagen-philharmonic"
        },
        offers: {
            price: "280",
            priceCurrency: "DKK",
            availability: "https://schema.org/InStock",
            url: "https://example.com/tickets/classical-symphony"
        },
        organizer: {
            name: "DR SymfoniOrkestret",
            url: "https://www.dr.dk/symfoni"
        }
    }
];

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('da-DK', {
        weekday: 'short',
        day: 'numeric',
        month: 'short',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}

function formatPrice(price, currency) {
    return new Intl.NumberFormat('da-DK', {
        style: 'currency',
        currency: currency
    }).format(price);
}

function createEventCard(event) {
    const locationIcon = '<svg class="inline w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd"></path></svg>';
    return 
        '<article class="event-card">' +
        '<div class="event-date">' + formatDate(event.startDate) + '</div>' +
        '<h2 class="event-title">' + event.name + '</h2>' +
        '<div class="event-venue">' + locationIcon + event.location.name + ' - ' + event.location.address.addressLocality + '</div>' +
        '<p class="event-description">' + event.description + '</p>' +
        '<div class="flex justify-between items-center mt-4">' +
        '<span class="text-primary font-medium">' + event.performer.name + '</span>' +
        '<span class="event-price">' + formatPrice(event.offers.price, event.offers.priceCurrency) + '</span>' +
        '</div>' +
        '</article>';
}

function renderEvents(events) {
    const container = document.getElementById('events-grid');
    container.innerHTML = events.map(createEventCard).join('');
}

function showEvents() {
    document.getElementById('skeleton-loader').classList.add('hidden');
    document.getElementById('events-container').classList.remove('hidden');
}

function init() {
    setTimeout(() => {
        renderEvents(sampleEvents);
        showEvents();
    }, 1500);
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    init();
}