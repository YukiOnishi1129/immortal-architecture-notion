export interface TemplateField {
  id: string;
  label: string;
  order: number;
  isRequired: boolean;
}

export interface TemplateOwner {
  id: string;
  firstName: string;
  lastName: string;
  thumbnail: string | null;
}

export interface Template {
  id: string;
  name: string;
  ownerId?: string;
  owner?: TemplateOwner;
  fields: TemplateField[];
  isUsed?: boolean;
  createdAt?: string;
  updatedAt: string;
}

export interface TemplateFilters {
  q?: string;
  page?: number;
  ownerId?: string;
  onlyMyTemplates?: boolean;
}
