export const goToForm = () => {
  const anchor = document.querySelector('#contact-form-anchor');

  if (anchor) {
    anchor.scrollIntoView({
      behavior: 'smooth',
      block: 'center',
    });
  }
}
