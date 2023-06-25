package mysql

import (
	"food_delivery_api/cfg"
	"food_delivery_api/pkg/model"
	"time"
)

func (s *Storage) ReadReportDashboard() (model.Dashboard, error) {
	var obj model.Dashboard
	var err error

	// forming summary
	var ttlCtg, ttlPrd, ttlTrx, ttlUom, ttlCst int64
	err = s.db.Model(&model.Category{}).Count(&ttlCtg).Error
	err = s.db.Model(&model.Product{}).Count(&ttlPrd).Error
	err = s.db.Model(&model.Transaction{}).Count(&ttlTrx).Error
	err = s.db.Model(&model.UOM{}).Count(&ttlUom).Error
	err = s.db.Model(&model.Transaction{}).Distinct("customer").Count(&ttlCst).Error
	if err != nil {
		return obj, err
	}

	sum := model.Summary{
		TotalCategory:    ttlCtg,
		TotalProduct:     ttlPrd,
		TotalTransaction: ttlTrx,
		TotalUOM:         ttlUom,
		TotalCustomer:    ttlCst,
	}

	// forming customer trx
	var ct []model.CustomerTrx
	s.db.Raw(`select customer, count(*) as total_trx, sum(total) as amount_trx, avg(total) as average_trx 
		from transactions group by customer`).
		Scan(&ct)

	// forming stock alert
	var sa []model.StockAlert
	s.db.Raw(`select p.id, p.code, p.name, p.qty, u.name as UOM from products p join uoms u on u.id = p.uom_id
		where qty < 100 order by qty;`).
		Scan(&sa)

	// forming top 10 trx
	var t10 []model.Top10Trx
	s.db.Raw(`select created_at, trx_id, created_by, customer, total from transactions order by total desc limit 10`).
		Scan(&t10)

	for i, v := range t10 {
		t10[i].TrxDate = v.CreatedAt.Format(cfg.AppTLayout)
	}

	// forming top 5 product
	var t5p []model.Top5Product
	s.db.Raw(`select product_id as id, name, price, sum(qty) as total_qty, sum(sub_total) / sum(qty) as average_amount,
    	sum(sub_total) as total_amount from transaction_lines group by product_id, name, price limit 5`).
		Scan(&t5p)

	obj = model.Dashboard{
		Summary:     sum,
		CustomerTrx: ct,
		StockAlert:  sa,
		Top10Trx:    t10,
		Top5Product: t5p,
	}

	return obj, nil
}

func (s *Storage) ReadReportChart() (model.Chart, error) {
	var obj model.Chart

	// forming daily trx chart
	var dr []model.DailyRow
	s.db.Raw(`select name as product, sum(qty) as qty, sum(sub_total) as amount from transaction_lines
		where DATE(created_at) = CURDATE() group by name`).
		Scan(&dr)

	var amount []model.Dataset
	var product []string
	var qtySold []int
	for _, v := range dr {
		amount = append(amount, model.Dataset{Label: v.Product, Data: []int64{v.Amount}})
		product = append(product, v.Product)
		qtySold = append(qtySold, v.Qty)
	}
	qty := []model.Dataset{{Label: "# of Qty", Data: qtySold}}

	dtac := model.ChartData{ChartType: "Vertical Bar Chart", Labels: []string{"Today"}, Datasets: amount}
	dtqc := model.ChartData{ChartType: "Doughnut Chart", Labels: product, Datasets: qty}

	// // forming monthly trx chart
	// var mr []model.MonthlyRow
	// s.db.Raw(`select month(created_at) as month, year(created_at) as year, sum(total) as amount from transactions group by year(created_at), month(created_at)`).
	// 	Scan(&mr)
	//
	// var mtc model.MonthlyTrxChart
	// mtc.ChartType = "Vertical Bar Chart"
	//
	// for i, v := range mr {
	// 	mr[i].Month = convertMonth(v.Month)
	// }

	// // forming multi type chart
	// var mtr []model.MultiTypeRow
	// var mtrLabel []string
	// var mtrCategory []string
	//
	// mChart := map[int]map[string]map[string]float64{}
	//
	// s.db.Raw(`select code from transaction_lines group by code`).Scan(&mtrCategory)
	// s.db.Raw(`select date_format(transactions.created_at, '%Y-%m') as month, transaction_lines.code as category,
	//    	sum(transaction_lines.sub_total) as total
	// 	from transactions
	//     join transaction_lines on transactions.id = transaction_lines.transaction_id
	// 	group by month, category`).
	// 	Scan(&mtr)
	//
	// for i, v := range mtr {
	// 	mtrLabel = append(mtrLabel, v.Month)
	// 	mChart[i] = map[string]map[string]float64{}
	// 	mChart[i][v.Month] = map[string]float64{}
	// 	mChart[i][v.Month][v.Category] = v.Total
	// }
	//
	// mtc := model.MultiTypeChart{
	// 	Labels:   nil,
	// 	Datasets: nil,
	// }

	obj = model.Chart{
		DailyTrxAmountChart: dtac,
		DailyTrxQtyChart:    dtqc,
		// MonthlyTrxChart:     mtc,
		// MultiTypeChart:      mtc,
	}

	return obj, nil
}

func convertMonth(strMonth string) string {
	date, _ := time.Parse("1", strMonth)
	return date.Month().String()
}
