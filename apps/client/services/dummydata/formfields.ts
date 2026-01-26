import { FormField } from '@/types/form.types'

export const FormFields: FormField[] = [
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: '06dbb551-3012-4b09-9c95-6f62fdd3cbb3',
    fieldLabel: 'short field',
    fieldType: 'text',
    isRequired: true,
    ordering: 0,
    options: []
  },
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: '5ead8b5f-776d-4937-8517-2e0e015a0b24',
    fieldLabel: 'New Field',
    fieldType: 'textArea',
    isRequired: true,
    ordering: 1,
    options: []
  },
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: 'da2ad7f4-4892-403f-85be-6e3690d7f478',
    fieldLabel: 'Number',
    fieldType: 'number',
    isRequired: true,
    ordering: 2,
    options: []
  },
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: '770dd847-c6f0-45a2-ad01-510f33bbcd86',
    fieldLabel: 'Email',
    fieldType: 'email',
    isRequired: false,
    ordering: 3,
    options: []
  },
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: 'baab330b-f80b-46a8-81dc-1bc406665c2b',
    fieldLabel: 'Phone',
    fieldType: 'phone',
    isRequired: true,
    ordering: 4,
    options: []
  },
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: '19e8956e-30fe-480a-8764-4a21f3294fef',
    fieldLabel: 'Url',
    fieldType: 'url',
    isRequired: true,
    ordering: 5,
    options: []
  },
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: '644d2a56-5725-4690-9e01-e2a68427b848',
    fieldLabel: 'Date',
    fieldType: 'date',
    isRequired: true,
    ordering: 6,
    options: []
  },
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: '4550a5a2-c5c7-435e-847e-93e2a0717acd',
    fieldLabel: 'Time',
    fieldType: 'time',
    isRequired: true,
    ordering: 7,
    options: []
  },
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: 'f344bb4b-b639-4763-98b5-1ea15937e218',
    fieldLabel: 'Date Time',
    fieldType: 'datetime',
    isRequired: false,
    ordering: 8,
    options: []
  },
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: '170579af-14db-4dac-8a89-65d24540843b',
    fieldLabel: 'Check Box',
    fieldType: 'checkbox',
    isRequired: true,
    ordering: 9,
    options: [
      {
        optionId: 'd71b927d-fdd4-4fcc-bab2-1a14ac5ca83e',
        fieldId: '170579af-14db-4dac-8a89-65d24540843b',
        optionLabel: 'Option 1',
        ordering: 0,
        isAnswer: true
      },
      {
        optionId: '687debd2-7f4c-468b-8d8b-27e159c42831',
        fieldId: '170579af-14db-4dac-8a89-65d24540843b',
        optionLabel: 'Option 3',
        ordering: 1,
        isAnswer: false
      },
      {
        optionId: '0d39e831-40d7-486d-aa3a-35bb5ca727e4',
        fieldId: '170579af-14db-4dac-8a89-65d24540843b',
        optionLabel: 'Option 2',
        ordering: 2,
        isAnswer: false
      }
    ]
  },
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: '0d0adc82-935e-4ce1-95bc-0d697a502b4c',
    fieldLabel: 'Single Choice',
    fieldType: 'radio',
    isRequired: true,
    ordering: 10,
    options: [
      {
        optionId: 'f4d21d73-35ac-400f-be93-27a1266d643c',
        fieldId: '0d0adc82-935e-4ce1-95bc-0d697a502b4c',
        optionLabel: 'option 1',
        ordering: 0,
        isAnswer: true
      },
      {
        optionId: 'bc900e0f-2632-42c3-9c90-d76081f320fe',
        fieldId: '0d0adc82-935e-4ce1-95bc-0d697a502b4c',
        optionLabel: 'Option 3',
        ordering: 1,
        isAnswer: false
      },
      {
        optionId: '3ec0b29a-9789-4d63-ae89-8a5a3d5db457',
        fieldId: '0d0adc82-935e-4ce1-95bc-0d697a502b4c',
        optionLabel: 'Option 2',
        ordering: 2,
        isAnswer: false
      }
    ]
  },
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: 'f07e7f8c-f231-4d35-8809-3a12cd7518df',
    fieldLabel: 'Select',
    fieldType: 'dropdown',
    isRequired: true,
    ordering: 11,
    options: [
      {
        optionId: '943f69fc-d289-4d5b-a7df-0a9fde3d3758',
        fieldId: 'f07e7f8c-f231-4d35-8809-3a12cd7518df',
        optionLabel: 'Option 1',
        ordering: 0,
        isAnswer: true
      },
      {
        optionId: '250656aa-9858-4a8a-b8e3-18639663cf25',
        fieldId: 'f07e7f8c-f231-4d35-8809-3a12cd7518df',
        optionLabel: 'Option 3',
        ordering: 1,
        isAnswer: false
      },
      {
        optionId: '43198993-014d-4361-9314-babd99f35140',
        fieldId: 'f07e7f8c-f231-4d35-8809-3a12cd7518df',
        optionLabel: 'Option 2',
        ordering: 2,
        isAnswer: false
      }
    ]
  },
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: '108b283f-6bd2-468f-bac7-286dd85070cc',
    fieldLabel: 'Multi Select',
    fieldType: 'multiselect',
    isRequired: true,
    ordering: 12,
    options: [
      {
        optionId: '0333c059-7069-43d4-bf24-683bc8ea11f4',
        fieldId: '108b283f-6bd2-468f-bac7-286dd85070cc',
        optionLabel: 'Option 1',
        ordering: 0,
        isAnswer: false
      },
      {
        optionId: '1cbaf587-fb6b-4f07-8573-403b2e5c40da',
        fieldId: '108b283f-6bd2-468f-bac7-286dd85070cc',
        optionLabel: 'Option 3',
        ordering: 1,
        isAnswer: false
      },
      {
        optionId: 'dd0721e8-e3d0-4a2d-8412-b0d57f0a1fd4',
        fieldId: '108b283f-6bd2-468f-bac7-286dd85070cc',
        optionLabel: 'Option 2',
        ordering: 2,
        isAnswer: false
      }
    ]
  },
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: '4070bf7d-d775-4869-b9b9-6fa31e5f8f04',
    fieldLabel: 'File Upload',
    fieldType: 'file',
    isRequired: false,
    ordering: 13,
    options: []
  },
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: '57ea040b-8101-4378-b413-ea68bf5bd346',
    fieldLabel: 'Image Upload',
    fieldType: 'image',
    isRequired: false,
    ordering: 14,
    options: []
  },
  {
    formId: 'a4ecc845-2a23-47c9-971e-e99150dc4b33',
    fieldId: '6cceed20-c723-47a7-8084-b9341037ee49',
    fieldLabel: 'Rating',
    fieldType: 'rating',
    isRequired: false,
    ordering: 15,
    options: []
  }
]
