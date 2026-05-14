import { render, screen, fireEvent } from '@testing-library/react'
import { MultiSelect } from './MultiSelect'
import { describe, it, expect, vi } from 'vitest'

const mockOptions = [
  { value: '1', label: 'Option 1' },
  { value: '2', label: 'Option 2' },
  { value: '3', label: 'Option 3' },
]

describe('MultiSelect', () => {
  it('renders placeholder when no value is selected', () => {
    render(<MultiSelect options={mockOptions} placeholder="Select options" onChange={() => {}} />)
    expect(screen.getByText('Select options')).toBeDefined()
  })

  it('renders selected options as badges', () => {
    render(<MultiSelect options={mockOptions} placeholder="Select options" value={['1', '2']} onChange={() => {}} />)
    expect(screen.getByText('Option 1')).toBeDefined()
    expect(screen.getByText('Option 2')).toBeDefined()
  })

  it('calls onChange when an option is toggled', () => {
    const handleChange = vi.fn()
    render(<MultiSelect options={mockOptions} placeholder="Select options" value={['1']} onChange={handleChange} />)
    
    // The SelectContent might not be visible in JSDOM unless triggered.
    // However, in MultiSelect.tsx, the options are always rendered inside SelectContent.
    
    fireEvent.click(screen.getByRole('combobox'))
    
    const option2 = screen.getByText('Option 2')
    fireEvent.click(option2)

    expect(handleChange).toHaveBeenCalledWith(['1', '2'])
  })

  it('removes an option when clicking its badge X icon', () => {
    const handleChange = vi.fn()
    render(<MultiSelect options={mockOptions} placeholder="Select options" value={['1', '2']} onChange={handleChange} />)
    
    // Find the X icon within the "Option 1" badge
    const badge1 = screen.getByText('Option 1').parentElement
    const xIcon = badge1?.querySelector('svg.lucide-x')
    
    if (xIcon) {
      fireEvent.click(xIcon)
    }

    expect(handleChange).toHaveBeenCalledWith(['2'])
  })
})
