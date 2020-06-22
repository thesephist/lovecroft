const LIST_NAME = window.location.pathname.split('/').reverse()[0];

const showActives = document.querySelector('.showActives');
const showAll = document.querySelector('.showAll');

const listSection = document.querySelector('.subscriberListSection');

showActives.addEventListener('click', () => {
    listSection.classList.remove('showingAll');
    listSection.classList.add('showingActives');
});

showAll.addEventListener('click', () => {
    listSection.classList.remove('showingActives');
    listSection.classList.add('showingAll');
});

const givenName = document.getElementById('add-sub-given-name');
const familyName = document.getElementById('add-sub-family-name');
const email = document.getElementById('add-sub-email');
const addForm = document.querySelector('.addForm');

function validate() {
    let valid = true;
    for (const el of [givenName, familyName, email]) {
        if (!el.value.trim()) {
            valid = false;
            el.classList.add('invalid');
        } else {
            el.classList.remove('invalid');
        }
    }

    if (!email.value.includes('@')) {
        valid = false;
        email.classList.add('invalid');
    } else {
        email.classList.remove('invalid');
    }

    return valid;
}

addForm.addEventListener('submit', evt => {
    evt.preventDefault();

    if (validate()) {
        fetch(`/subscribe/${LIST_NAME}`, {
            method: 'POST',
            body: JSON.stringify({
                givenName: givenName.value.trim(),
                familyName: familyName.value.trim(),
                email: email.value.trim(),
            }),
        }).then(resp => {
            if (resp.status === 500) {
                resp.text().then(t => alert(t));
            } else {
                window.location.reload();
            }
        }).catch(e => {
            alert(e);
            console.error(e);
        });
    }
});