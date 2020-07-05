const el = {
    from: document.getElementById('issue-from'),
    list: document.getElementById('issue-list'),
    subject: document.getElementById('issue-subject'),
    body: document.getElementById('issue-body'),
}

const sendForm = document.querySelector('.sendForm');
const sendButton = document.getElementById('sendButton');

function validate() {
    let valid = true;
    for (const input of Object.values(el)) {
        if (!input.value.trim()) {
            input.classList.add('invalid');
            valid = false;
        } else {
            input.classList.remove('invalid');
        }
    }

    return valid;
}

sendForm.addEventListener('submit', evt => {
    evt.preventDefault();

    if (validate()) {
        fetch(`/admin/send/${el.list.value}`, {
            method: 'POST',
            body: JSON.stringify({
                from: el.from.value.trim(),
                subject: el.subject.value.trim(),
                body: el.body.value.trim(),
            }),
        }).then(resp => {
            if (resp.status === 500 || resp.status === 409) {
                resp.text().then(t => alert(t));
            } else {
                alert('Queued successfully!');
            }
        }).catch(e => {
            alert(e);
            console.error(e);
        });
    }
});