const listName = document.getElementById('add-sub-name');
const addForm = document.querySelector('.addForm');

function validate() {
    if (!listName.value.trim() || listName.value.includes('/')) {
        listName.classList.add('invalid');
        return false;
    } else {
        listName.classList.remove('invalid');
        return true;
    }
}

addForm.addEventListener('submit', evt => {
    evt.preventDefault();

    if (validate()) {
        fetch(`/admin/create-list/${listName.value.trim()}`, {
            method: 'POST',
        }).then(resp => {
            if (resp.status === 500 || resp.status === 409) {
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