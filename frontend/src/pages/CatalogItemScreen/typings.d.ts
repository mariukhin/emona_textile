type CatalogItemData = {
  id: number;
  title: string;
  buttonText?: string;
  imageUrl: string;
  items: PageItems[];
};

type PageItems = {
  id: number;
  title: string;
  imageUrl: string;
  description: string[];
}
