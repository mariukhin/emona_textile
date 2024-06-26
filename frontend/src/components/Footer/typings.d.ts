type FooterData = {
  title: string;
  subItems: FooterSubItem[]
};

type FooterSubItem = {
  label: string;
  href?: string;
  onClick?: MouseEventHandler<HTMLAnchorElement> & MouseEventHandler<HTMLSpanElement>;
};
